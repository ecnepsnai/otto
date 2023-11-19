package server

import (
	"fmt"
	"sort"

	"github.com/ecnepsnai/web"
)

func (h *handle) RunbookList(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	runbooks := RunbookStore.AllRunbooks()
	sort.Slice(runbooks, func(i int, j int) bool {
		return runbooks[i].Name < runbooks[j].Name
	})

	return runbooks, nil, nil
}

func (h *handle) RunbookGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	runbook := RunbookStore.RunbookWithID(id)
	if runbook == nil {
		return nil, nil, web.ValidationError("No runbook with ID %s", id)
	}

	return runbook, nil, nil
}

func (h *handle) RunbookNew(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	if !session.User().Permissions.CanModifyRunbooks {
		EventStore.UserPermissionDenied(session.User().Username, "Create new runbook")
		return nil, nil, web.ValidationError("Permission denied")
	}

	params := Runbook{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	runbook, err := RunbookStore.New(params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.RunbookAdded(runbook, session.Username)

	return runbook, nil, nil
}

func (h *handle) RunbookEdit(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifyRunbooks {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Modify runbook %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	if RunbookStore.RunbookWithID(id) == nil {
		return nil, nil, web.ValidationError("No runbook with ID %s", id)
	}

	runbook := Runbook{}
	if err := request.DecodeJSON(&runbook); err != nil {
		return nil, nil, err
	}

	newRunbook, err := RunbookStore.Edit(id, runbook)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.RunbookModified(newRunbook, session.Username)

	return newRunbook, nil, nil
}

func (h *handle) RunbookDelete(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifyRunbooks {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Delete runbook %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	if RunbookStore.RunbookWithID(id) == nil {
		return nil, nil, web.ValidationError("No runbook with ID %s", id)
	}

	runbook, err := RunbookStore.Delete(id)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.RunbookDeleted(runbook, session.Username)

	return true, nil, nil
}
