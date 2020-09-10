package server

import (
	"github.com/ecnepsnai/web"
)

func (h *handle) GroupList(request web.Request) (interface{}, *web.Error) {
	groups, err := GroupStore.AllGroups()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return groups, nil
}

func (h *handle) GroupGetMembership(request web.Request) (interface{}, *web.Error) {
	return GetGroupCache(), nil
}

func (h *handle) GroupGet(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	group, err := GroupStore.GroupWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	return group, nil
}

func (h *handle) GroupGetHosts(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	group, err := GroupStore.GroupWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	hosts, err := group.Hosts()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return hosts, nil
}

func (h *handle) GroupSetHosts(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	type params struct {
		Hosts []string
	}

	r := params{}
	if err := request.Decode(&r); err != nil {
		return nil, err
	}

	group, err := GroupStore.GroupWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	for _, hostID := range r.Hosts {
		host, err := HostStore.HostWithID(hostID)
		if err != nil {
			if err.Server {
				return nil, web.CommonErrors.ServerError
			}
			return nil, web.ValidationError(err.Message)
		}
		if StringSliceContains(id, host.GroupIDs) {
			continue
		}

		if _, err := HostStore.EditHost(host, editHostParameters{
			Name:        host.Name,
			Address:     host.Address,
			Port:        host.Port,
			PSK:         host.PSK,
			Enabled:     host.Enabled,
			GroupIDs:    append(host.GroupIDs, id),
			Environment: host.Environment,
		}); err != nil {
			return nil, web.CommonErrors.ServerError
		}
	}

	hosts, err := group.Hosts()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return hosts, nil
}

func (h *handle) GroupGetScripts(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	group, err := GroupStore.GroupWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	scripts, err := group.Scripts()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return scripts, nil
}

func (h *handle) GroupGetSchedules(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedules, err := ScheduleStore.AllSchedulesForGroup(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return schedules, nil
}

func (h *handle) GroupNew(request web.Request) (interface{}, *web.Error) {
	params := newGroupParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	group, err := GroupStore.NewGroup(params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return group, nil
}

func (h *handle) GroupEdit(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	group, err := GroupStore.GroupWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	params := editGroupParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	group, err = GroupStore.EditGroup(group, params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return group, nil
}

func (h *handle) GroupDelete(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	group, err := GroupStore.GroupWithID(id)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	if err := GroupStore.DeleteGroup(group); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return true, nil
}
