package server

import (
	"fmt"
	"time"

	"github.com/ecnepsnai/web"
)

func (h *handle) RequestNew(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	type requestParams struct {
		HostID   string
		Action   string
		ScriptID string
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
		script := ScriptStore.ScriptWithID(r.ScriptID)
		if script == nil {
			return nil, nil, web.ValidationError("No script with ID %s", r.ScriptID)
		}

		result, err := host.RunScript(script, nil, nil)
		if err != nil {
			return nil, nil, web.CommonErrors.ServerError
		}

		EventStore.ScriptRun(script, host, &result.Result, nil, session.Username)

		return result, nil, nil
	}

	return nil, nil, nil
}

func (h handle) RequestStream(request web.Request, conn *web.WSConn) {
	session := request.UserData.(*Session)
	defer conn.Close()

	type requestParams struct {
		HostID   string
		Action   string
		ScriptID string
	}
	type requestResponse struct {
		Code   int           `json:"Code,omitempty"`
		Error  string        `json:"Error,omitempty"`
		Stdout string        `json:"Stdout,omitempty"`
		Stderr string        `json:"Stderr,omitempty"`
		Result *ScriptResult `json:"Result,omitempty"`
	}

	writeMessage := func(m requestResponse) {
		log.Debug("ws send %d", m.Code)
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
		script := ScriptStore.ScriptWithID(r.ScriptID)
		if script == nil {
			writeMessage(requestResponse{
				Code:  RequestResponseCodeError,
				Error: fmt.Sprintf("No script with ID %s", r.ScriptID),
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

		running := true
		cancel := make(chan bool)
		go func() {
			type cancelParams struct {
				Cancel bool
			}
			for running {
				cancelRequest := cancelParams{}
				if err := conn.ReadJSON(&cancelRequest); err != nil {
					log.PError("Error reading from websocket connection", map[string]interface{}{
						"error": err.Error(),
					})
					running = false
					break
				}
				if cancelRequest.Cancel {
					log.PWarn("Request to cancel running script", map[string]interface{}{
						"script_id": r.ScriptID,
						"host_id":   r.HostID,
					})
					cancel <- true
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()
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
		}, cancel)
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
			Code:   RequestResponseCodeFinished,
			Result: result,
		})
		running = false
		return
	}
}
