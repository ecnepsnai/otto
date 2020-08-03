package server

import (
	"strings"

	"github.com/ecnepsnai/web"
)

func (h *handle) OptionsGet(request web.Request) (interface{}, *web.Error) {
	return Options, nil
}

func (h *handle) OptionsSet(request web.Request) (interface{}, *web.Error) {
	options := OttoOptions{}

	if err := request.Decode(&options); err != nil {
		return nil, web.CommonErrors.BadRequest
	}

	if !strings.HasPrefix(options.General.ServerURL, "http") {
		return nil, web.ValidationError("Server URL must include protocol")
	}

	if !strings.HasSuffix(options.General.ServerURL, "/") {
		options.General.ServerURL = options.General.ServerURL + "/"
	}

	if err := options.Save(); err != nil {
		return nil, web.CommonErrors.ServerError
	}

	return options, nil
}
