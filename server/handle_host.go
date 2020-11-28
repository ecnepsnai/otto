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

func (h *handle) HostGetGroups(request web.Request) (interface{}, *web.Error) {
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

	groups, err := host.Groups()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return groups, nil
}

func (h *handle) HostGetSchedules(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedules, err := ScheduleStore.AllSchedulesForHost(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return schedules, nil
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
	session := request.UserData.(*Session)

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

	EventStore.HostAdded(host, session.Username)

	return host, nil
}

func (h *handle) HostEdit(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

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

	EventStore.HostModified(host, session.Username)

	return host, nil
}

func (h *handle) HostDelete(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

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

	EventStore.HostDeleted(host, session.Username)

	return true, nil
}
