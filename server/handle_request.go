package server

import (
	"fmt"
	"time"

	"github.com/ecnepsnai/web"
)

func (h *handle) RequestNew(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	type requestParams struct {
		HostID   string
		Action   string
		ScriptID string
	}

	r := requestParams{}
	if err := request.DecodeJSON(&r); err != nil {
		return nil, err
	}

	host := HostCache.ByID(r.HostID)
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", r.HostID)
	}

	if !IsClientAction(r.Action) {
		return nil, web.ValidationError("Unknown action %s", r.Action)
	}

	if r.Action == ClientActionPing {
		if err := host.Ping(); err != nil {
			return false, nil
		}
		return true, nil
	} else if r.Action == ClientActionRunScript {
		script := ScriptStore.ScriptWithID(r.ScriptID)
		if script == nil {
			return nil, web.ValidationError("No script with ID %s", r.ScriptID)
		}

		result, err := host.RunScript(script, nil, nil)
		if err != nil {
			return nil, web.CommonErrors.ServerError
		}

		EventStore.ScriptRun(script, host, &result.Result, nil, session.Username)

		return result, nil
	} else if r.Action == ClientActionExitClient {
		if err := host.ExitClient(); err != nil {
			return false, nil
		}
		return true, nil
	}

	return nil, nil
}

func (h handle) RequestStream(request web.Request, conn web.WSConn) {
	const (
		requestResponseCodeOutput    = 100
		requestResponseCodeKeepalive = 101
		requestResponseCodeError     = 400
		requestResponseCodeFinished  = 200
	)

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

	r := requestParams{}
	if err := conn.ReadJSON(&r); err != nil {
		conn.WriteJSON(requestResponse{
			Code:  requestResponseCodeError,
			Error: "Invalid request",
		})
		return
	}

	host := HostCache.ByID(r.HostID)
	if host == nil {
		conn.WriteJSON(requestResponse{
			Code:  requestResponseCodeError,
			Error: fmt.Sprintf("No host with ID %s", r.HostID),
		})
		return
	}

	if !IsClientAction(r.Action) {
		conn.WriteJSON(requestResponse{
			Code:  requestResponseCodeError,
			Error: fmt.Sprintf("Unknown action %s", r.Action),
		})
		return
	}

	if r.Action == ClientActionPing {
		if err := host.Ping(); err != nil {
			conn.WriteJSON(requestResponse{
				Code:  requestResponseCodeError,
				Error: fmt.Sprintf("Error pinging host %s", err.Error),
			})
			return
		}
		conn.WriteJSON(requestResponse{Code: 200})
		return
	} else if r.Action == ClientActionRunScript {
		script := ScriptStore.ScriptWithID(r.ScriptID)
		if script == nil {
			conn.WriteJSON(requestResponse{
				Code:  requestResponseCodeError,
				Error: fmt.Sprintf("No script with ID %s", r.ScriptID),
			})
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
				conn.ReadJSON(&cancelRequest)
				if cancelRequest.Cancel {
					log.Warn("Request to cancel running script '%s' on host '%s'", r.ScriptID, r.HostID)
					cancel <- true
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()
		go func() {
			lastKA := time.Now().AddDate(0, 0, -1)
			for running {
				if time.Since(lastKA) > 10*time.Second {
					conn.WriteJSON(requestResponse{
						Code: requestResponseCodeKeepalive,
					})
					lastKA = time.Now()
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()

		result, err := host.RunScript(script, func(stdout, stderr []byte) {
			conn.WriteJSON(requestResponse{
				Code:   requestResponseCodeOutput,
				Stdout: string(stdout),
				Stderr: string(stderr),
			})
		}, cancel)
		if err != nil {
			conn.WriteJSON(requestResponse{
				Code:  requestResponseCodeError,
				Error: err.Message,
			})
			running = false
			return
		}

		EventStore.ScriptRun(script, host, &result.Result, nil, session.Username)
		conn.WriteJSON(requestResponse{
			Code:   requestResponseCodeFinished,
			Result: result,
		})
		running = false
		return
	} else if r.Action == ClientActionExitClient {
		if err := host.ExitClient(); err != nil {
			conn.WriteJSON(requestResponse{Code: requestResponseCodeFinished})
			return
		}
		conn.WriteJSON(requestResponse{Code: requestResponseCodeFinished})
		return
	}
}
