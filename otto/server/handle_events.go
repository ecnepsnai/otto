package server

import (
	"strconv"

	"github.com/ecnepsnai/web"
)

func (h *handle) EventsGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	if !session.User().Permissions.CanAccessAuditLog {
		EventStore.UserPermissionDenied(session.User().Username, "Access audit log")
		return nil, nil, web.ValidationError("Permission denied")
	}

	if len(request.HTTP.URL.Query()["c"]) < 1 {
		return nil, nil, web.ValidationError("Must specify max number of events")
	}
	count, cerr := strconv.Atoi(request.HTTP.URL.Query()["c"][0])
	if cerr != nil {
		return nil, nil, web.ValidationError("Invalid count")
	}

	events, err := EventStore.LastEvents(count)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	return events, nil, nil
}
