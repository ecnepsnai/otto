package server

import (
	"fmt"
	"time"

	"github.com/ecnepsnai/web"
)

func (h *handle) RequestNew(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	type requestParams struct {
		HostID     string
		Action     string
		ResourceID string
	}

	r := requestParams{}
	if err := request.DecodeJSON(&r); err != nil {
		return nil, nil, err
	}

	host := HostCache.ByID(r.HostID)
	if host == nil {
		return nil, nil, web.ValidationError("No host with ID %s", r.HostID)
	}

	if !IsAgentAction(r.Action) {
		return nil, nil, web.ValidationError("Unknown action %s", r.Action)
	}

	if r.Action == AgentActionPing {
		if err := host.Ping(); err != nil {
			return false, nil, nil
		}
		return true, nil, nil
	} else if r.Action == AgentActionRunScript {
		script := ScriptStore.ScriptWithID(r.ResourceID)
		if script == nil {
			return nil, nil, web.ValidationError("No script with ID %s", r.ResourceID)
		}

		if session.User().Permissions.ScriptRunLevel < script.RunLevel {
			EventStore.UserPermissionDenied(session.Username, fmt.Sprintf("attempt to cancel script with higher run level: %s", script.Name))
			return nil, nil, web.CommonErrors.Forbidden
		}

		result, err := host.RunScript(script, nil)
		if err != nil {
			return nil, nil, web.CommonErrors.ServerError
		}

		EventStore.ScriptRun(script, host, &result.Result, nil, session.Username)

		return result, nil, nil
	} else if r.Action == AgentActionStartRunbook {
		runbook := RunbookStore.RunbookWithID(r.ResourceID)
		if runbook == nil {
			return nil, nil, web.ValidationError("No runbook with ID %s", r.ResourceID)
		}

		if session.User().Permissions.ScriptRunLevel < runbook.RunLevel {
			EventStore.UserPermissionDenied(session.Username, fmt.Sprintf("attempt to cancel runbook with higher run level: %s", runbook.Name))
			return nil, nil, web.CommonErrors.Forbidden
		}

		var results []ScriptResult

		for _, scriptId := range runbook.ScriptIDs {
			script := ScriptStore.ScriptWithID(scriptId)
			if script == nil {
				return nil, nil, web.ValidationError("No script with ID %s", r.ResourceID)
			}

			result, err := host.RunScript(script, nil)
			if err != nil {
				return nil, nil, web.CommonErrors.ServerError
			}

			EventStore.ScriptRun(script, host, &result.Result, nil, session.Username)

			results = append(results, *result)
		}

		return results, nil, nil
	}

	return nil, nil, nil
}

func (h *handle) RequestCancel(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	type tCancelRequest struct {
		HostID   string
		ScriptID string
	}
	cancelRequest := tCancelRequest{}
	if err := request.DecodeJSON(&cancelRequest); err != nil {
		return nil, nil, err
	}

	script := ScriptCache.ByID(cancelRequest.ScriptID)
	if script == nil {
		return nil, nil, web.ValidationError("No script found with id %s", cancelRequest.ScriptID)
	}

	host := HostCache.ByID(cancelRequest.HostID)
	if host == nil {
		return nil, nil, web.ValidationError("No host found with id %s", cancelRequest.HostID)
	}

	if session.User().Permissions.ScriptRunLevel < script.RunLevel {
		EventStore.UserPermissionDenied(session.Username, fmt.Sprintf("attempt to cancel script with higher run level: %s", script.Name))
		return nil, nil, web.CommonErrors.Forbidden
	}

	if err := host.CancelScript(script.Name); err != nil {
		log.PError("Error cancelling script on host", map[string]interface{}{
			"host_id":     host.ID,
			"script_name": script.Name,
			"error":       err.Error(),
		})
		return nil, nil, web.ValidationError("%s", err.Error())
	}

	return true, nil, nil
}

func (h handle) RequestStream(request web.Request, conn *web.WSConn) {
	session := request.UserData.(*Session)
	defer conn.Close()

	type requestParams struct {
		HostID     string
		Action     string
		ResourceID string
	}
	type requestResponse struct {
		Code     int           `json:"Code,omitempty"`
		Error    string        `json:"Error,omitempty"`
		Stdout   string        `json:"Stdout,omitempty"`
		Stderr   string        `json:"Stderr,omitempty"`
		Result   *ScriptResult `json:"Result,omitempty"`
		ScriptId string        `json:"ScriptId,omitempty"`
	}

	writeMessage := func(m requestResponse) {
		if err := conn.WriteJSON(m); err != nil {
			log.PError("Error sending websocket message", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	r := requestParams{}
	if err := conn.ReadJSON(&r); err != nil {
		writeMessage(requestResponse{
			Code:  RequestResponseCodeError,
			Error: "Invalid request",
		})
		return
	}

	host := HostCache.ByID(r.HostID)
	if host == nil {
		writeMessage(requestResponse{
			Code:  RequestResponseCodeError,
			Error: fmt.Sprintf("No host with ID %s", r.HostID),
		})
		return
	}

	if !IsAgentAction(r.Action) {
		writeMessage(requestResponse{
			Code:  RequestResponseCodeError,
			Error: fmt.Sprintf("Unknown action %s", r.Action),
		})
		return
	}

	if r.Action == AgentActionPing {
		if err := host.Ping(); err != nil {
			writeMessage(requestResponse{
				Code:  RequestResponseCodeError,
				Error: fmt.Sprintf("Error pinging host %s", err.Error()),
			})
			return
		}
		writeMessage(requestResponse{Code: 200})
		return
	} else if r.Action == AgentActionRunScript {
		script := ScriptStore.ScriptWithID(r.ResourceID)
		if script == nil {
			writeMessage(requestResponse{
				Code:  RequestResponseCodeError,
				Error: fmt.Sprintf("No script with ID %s", r.ResourceID),
			})
			return
		}

		if session.User().Permissions.ScriptRunLevel < script.RunLevel {
			writeMessage(requestResponse{
				Code:  RequestResponseCodeError,
				Error: "Permission denied",
			})
			EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Run script %s", script.ID))
			return
		}

		writeMessage(requestResponse{
			Code:     RequestResponseCodeStartScript,
			ScriptId: script.ID,
		})

		running := true
		go func() {
			lastKA := time.Now().AddDate(0, 0, -1)
			for running {
				if time.Since(lastKA) > 10*time.Second {
					writeMessage(requestResponse{
						Code: RequestResponseCodeKeepalive,
					})
					lastKA = time.Now()
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()

		result, err := host.RunScript(script, func(stdout, stderr []byte) {
			writeMessage(requestResponse{
				Code:   RequestResponseCodeOutput,
				Stdout: string(stdout),
				Stderr: string(stderr),
			})
		})
		if err != nil {
			writeMessage(requestResponse{
				Code:  RequestResponseCodeError,
				Error: err.Error(),
			})
			running = false
			return
		}

		EventStore.ScriptRun(script, host, &result.Result, nil, session.Username)
		writeMessage(requestResponse{
			Code:   RequestResponseCodeFinishedScript,
			Result: result,
		})
		running = false
		return
	}
}
