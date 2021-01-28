package server

import (
	"github.com/ecnepsnai/web"
)

func (h *handle) ScriptList(request web.Request) (interface{}, *web.Error) {
	return ScriptStore.AllScripts(), nil
}

func (h *handle) ScriptGet(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, web.ValidationError("No script with ID %s", id)
	}

	return script, nil
}

func (h *handle) ScriptGetGroups(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, web.ValidationError("No script with ID %s", id)
	}

	return script.Groups(), nil
}

func (h *handle) ScriptGetHosts(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, web.ValidationError("No script with ID %s", id)
	}

	return script.Hosts(), nil
}

func (h *handle) ScriptGetSchedules(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedules := ScheduleStore.AllSchedulesForScript(id)
	return schedules, nil
}

func (h *handle) ScriptGetAttachments(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, web.ValidationError("No script with ID %s", id)
	}

	files, err := script.Attachments()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return files, nil
}

func (h *handle) ScriptSetGroups(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	type params struct {
		Groups []string
	}

	r := params{}
	if err := request.Decode(&r); err != nil {
		return nil, err
	}

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, web.ValidationError("No script with ID %s", id)
	}

	if err := ScriptStore.SetGroups(script, r.Groups); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return script.Hosts(), nil
}

func (h *handle) ScriptNew(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	params := newScriptParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	script, err := ScriptStore.NewScript(params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.ScriptAdded(script, session.Username)

	return script, nil
}

func (h *handle) ScriptEdit(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Params.ByName("id")

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, web.ValidationError("No script with ID %s", id)
	}

	params := editScriptParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	script, err := ScriptStore.EditScript(script, params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.ScriptModified(script, session.Username)

	return script, nil
}

func (h *handle) ScriptDelete(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Params.ByName("id")

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, web.ValidationError("No script with ID %s", id)
	}

	if err := ScriptStore.DeleteScript(script); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.ScriptDeleted(script, session.Username)

	return true, nil
}
