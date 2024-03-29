package server

import (
	"fmt"
	"sort"

	"github.com/ecnepsnai/web"
)

func (h *handle) ScriptList(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	scripts := ScriptStore.AllScripts()
	sort.Slice(scripts, func(i int, j int) bool {
		return scripts[i].Name < scripts[j].Name
	})

	// Hide secret environment variables if the user cannot modify them
	if !session.User().Permissions.CanModifyHosts {
		for i, script := range scripts {
			for y, env := range script.Environment {
				if env.Secret {
					scripts[i].Environment[y].Value = ""
				}
			}
		}
	}

	return scripts, nil, nil
}

func (h *handle) ScriptGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]
	session := request.UserData.(*Session)

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, nil, web.ValidationError("No script with ID %s", id)
	}

	// Hide secret environment variables if the user cannot modify them
	if !session.User().Permissions.CanModifyHosts {
		for i, env := range script.Environment {
			if env.Secret {
				script.Environment[i].Value = ""
			}
		}
	}

	return script, nil, nil
}

func (h *handle) ScriptGetGroups(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, nil, web.ValidationError("No script with ID %s", id)
	}
	groups := script.Groups()
	sort.Slice(groups, func(i int, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups, nil, nil
}

func (h *handle) ScriptGetHosts(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, nil, web.ValidationError("No script with ID %s", id)
	}
	hosts := script.Hosts()
	sort.Slice(hosts, func(i int, j int) bool {
		return hosts[i].HostName < hosts[j].HostName
	})

	return hosts, nil, nil
}

func (h *handle) ScriptGetSchedules(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	schedules := ScheduleStore.AllSchedulesForScript(id)
	sort.Slice(schedules, func(i int, j int) bool {
		return schedules[i].Name < schedules[j].Name
	})

	return schedules, nil, nil
}

func (h *handle) ScriptGetAttachments(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, nil, web.ValidationError("No script with ID %s", id)
	}

	files, err := script.Attachments()
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}
	sort.Slice(files, func(i int, j int) bool {
		return files[i].Name < files[j].Name
	})

	return files, nil, nil
}

func (h *handle) ScriptSetGroups(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifyScripts {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Modify groups for script %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	type params struct {
		Groups []string
	}

	r := params{}
	if err := request.DecodeJSON(&r); err != nil {
		return nil, nil, err
	}

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, nil, web.ValidationError("No script with ID %s", id)
	}

	if err := ScriptStore.SetGroups(script, r.Groups); err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	return script.Hosts(), nil, nil
}

func (h *handle) ScriptNew(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	if !session.User().Permissions.CanModifyScripts {
		EventStore.UserPermissionDenied(session.User().Username, "Create new script")
		return nil, nil, web.ValidationError("Permission denied")
	}

	params := newScriptParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	script, err := ScriptStore.NewScript(params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.ScriptAdded(script, session.Username)

	return script, nil, nil
}

func (h *handle) ScriptEdit(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifyScripts {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Modify script %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, nil, web.ValidationError("No script with ID %s", id)
	}

	params := editScriptParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	script, err := ScriptStore.EditScript(script, params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.ScriptModified(script, session.Username)

	return script, nil, nil
}

func (h *handle) ScriptDelete(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifyScripts {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Delete script %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	script := ScriptStore.ScriptWithID(id)
	if script == nil {
		return nil, nil, web.ValidationError("No script with ID %s", id)
	}

	if err := ScriptStore.DeleteScript(script); err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.ScriptDeleted(script, session.Username)

	return true, nil, nil
}
