package server

import "github.com/ecnepsnai/web"

func (h *handle) RequestNew(request web.Request) (interface{}, *web.Error) {
	type requestParams struct {
		HostID   string
		Action   string
		ScriptID string
	}

	r := requestParams{}
	if err := request.Decode(&r); err != nil {
		return nil, err
	}

	host, err := HostStore.HostWithID(r.HostID)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
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
		script, err := ScriptStore.ScriptWithID(r.ScriptID)
		if err != nil {
			if err.Server {
				return nil, web.CommonErrors.ServerError
			}
			return nil, web.ValidationError(err.Message)
		}
		if script == nil {
			return nil, web.ValidationError("No script with ID %s", r.ScriptID)
		}

		result, err := host.RunScript(script)
		if err != nil {
			return nil, web.ValidationError(err.Message)
		}

		return result, nil
	} else if r.Action == ClientActionExitClient {
		if err := host.ExitClient(); err != nil {
			return false, nil
		}
		return true, nil
	}

	return nil, nil
}
