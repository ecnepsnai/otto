package server

import "github.com/ecnepsnai/web"

func (h *handle) SystemSearch(request web.Request) (interface{}, *web.Error) {
	type systemSearchRequest struct {
		Query string
	}

	req := systemSearchRequest{}
	if err := request.Decode(&req); err != nil {
		return nil, err
	}

	return SystemSearch(req.Query), nil
}
