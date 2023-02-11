package server

import "github.com/ecnepsnai/web"

func (h *handle) SystemSearch(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	type systemSearchRequest struct {
		Query string
	}

	req := systemSearchRequest{}
	if err := request.DecodeJSON(&req); err != nil {
		return nil, nil, err
	}

	return SystemSearch(req.Query), nil, nil
}
