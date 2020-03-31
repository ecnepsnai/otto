package server

import (
	"github.com/ecnepsnai/web"
)

func (h *handle) HostList(request web.Request) (interface{}, *web.Error) {
	hosts, err := HostStore.AllHosts()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return hosts, nil
}

func (h *handle) HostGet(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host, err := HostStore.HostWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	return host, nil
}

func (h *handle) HostGetScripts(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host, err := HostStore.HostWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	return host.Scripts(), nil
}

func (h *handle) HostNew(request web.Request) (interface{}, *web.Error) {
	params := newHostParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	host, err := HostStore.NewHost(params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return host, nil
}

func (h *handle) HostEdit(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host, err := HostStore.HostWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	params := editHostParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	host, err = HostStore.EditHost(host, params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return host, nil
}

func (h *handle) HostDelete(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host, err := HostStore.HostWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	if err := HostStore.DeleteHost(host); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return true, nil
}
