package server

import (
	"strings"

	"github.com/ecnepsnai/web"
)

func (h *handle) OptionsGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	if !session.User().Permissions.CanModifySystem {
		EventStore.UserPermissionDenied(session.User().Username, "Access system settings")
		return nil, nil, web.ValidationError("Permission denied")
	}

	return Options, nil, nil
}

func (h *handle) OptionsSet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	if !session.User().Permissions.CanModifySystem {
		EventStore.UserPermissionDenied(session.User().Username, "Modify system settings")
		return nil, nil, web.ValidationError("Permission denied")
	}

	options := OttoOptions{}

	if err := request.DecodeJSON(&options); err != nil {
		return nil, nil, web.CommonErrors.BadRequest
	}

	if !strings.HasPrefix(options.General.ServerURL, "http") {
		return nil, nil, web.ValidationError("Server URL must include protocol")
	}

	if !strings.HasSuffix(options.General.ServerURL, "/") {
		options.General.ServerURL = options.General.ServerURL + "/"
	}

	if err := options.Validate(); err != nil {
		return nil, nil, web.ValidationError(err.Error())
	}

	hash, didChange := options.Save()
	if didChange {
		EventStore.ServerOptionsModified(hash, session.Username)
	}

	return options, nil, nil
}
