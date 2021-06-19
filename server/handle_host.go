package server

import (
	"sort"

	"github.com/ecnepsnai/web"
)

func (h *handle) HostList(request web.Request) (interface{}, *web.Error) {
	hosts := HostStore.AllHosts()
	sort.Slice(hosts, func(i int, j int) bool {
		return hosts[i].Name < hosts[j].Name
	})

	return hosts, nil
}

func (h *handle) HostGet(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host := HostStore.HostWithID(id)
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	return host, nil
}

func (h *handle) HostGetGroups(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host := HostStore.HostWithID(id)
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

	sort.Slice(groups, func(i int, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups, nil
}

func (h *handle) HostGetSchedules(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedules := ScheduleStore.AllSchedulesForHost(id)
	sort.Slice(schedules, func(i int, j int) bool {
		return schedules[i].Name < schedules[j].Name
	})

	return schedules, nil
}

func (h *handle) HostRotatePSK(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host := HostStore.HostWithID(id)
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	newPSK, err := host.RotatePSKNow()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return newPSK, nil
}

func (h *handle) HostTriggerHeartbeat(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host := HostStore.HostWithID(id)
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	host.Ping()
	return heartbeatStore.LastHeartbeat(host), nil
}

func (h *handle) HostGetScripts(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	host := HostStore.HostWithID(id)
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	scripts := host.Scripts()
	sort.Slice(scripts, func(i int, j int) bool {
		return scripts[i].ScriptName < scripts[j].ScriptName
	})

	return scripts, nil
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

	host := HostStore.HostWithID(id)
	if host == nil {
		return nil, web.ValidationError("No host with ID %s", id)
	}

	params := editHostParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	host, err := HostStore.EditHost(host, params)
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

	host := HostStore.HostWithID(id)
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
