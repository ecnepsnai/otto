package server

import (
	"sort"

	"github.com/ecnepsnai/web"
)

func (h *handle) GroupList(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	groups := GroupStore.AllGroups()
	sort.Slice(groups, func(i int, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups, nil, nil
}

func (h *handle) GroupGetMembership(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	return GroupCache.Membership(), nil, nil
}

func (h *handle) GroupGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	group := GroupCache.ByID(id)
	if group == nil {
		return nil, nil, web.ValidationError("No group with ID %s", id)
	}

	return group, nil, nil
}

func (h *handle) GroupGetHosts(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	group := GroupCache.ByID(id)
	if group == nil {
		return nil, nil, web.ValidationError("No group with ID %s", id)
	}

	hosts, err := group.Hosts()
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}
	sort.Slice(hosts, func(i int, j int) bool {
		return hosts[i].Name < hosts[j].Name
	})

	return hosts, nil, nil
}

func (h *handle) GroupSetHosts(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	type params struct {
		Hosts []string
	}

	r := params{}
	if err := request.DecodeJSON(&r); err != nil {
		return nil, nil, err
	}

	group := GroupCache.ByID(id)
	if group == nil {
		return nil, nil, web.ValidationError("No group with ID %s", id)
	}

	var addedHosts []string
	var removedHosts []string

	currentHosts := group.HostIDs()
	for _, hostID := range r.Hosts {
		if !stringSliceContains(hostID, currentHosts) {
			addedHosts = append(addedHosts, hostID)
		}
	}
	for _, hostID := range currentHosts {
		if !stringSliceContains(hostID, r.Hosts) {
			removedHosts = append(removedHosts, hostID)
		}
	}

	log.Debug("Will add hosts to group %s: %+v", id, addedHosts)
	log.Debug("Will remove hosts from group %s: %+v", id, removedHosts)

	for _, hostID := range addedHosts {
		host := HostCache.ByID(hostID)
		if stringSliceContains(id, host.GroupIDs) {
			continue
		}

		if _, err := HostStore.EditHost(host, editHostParameters{
			Name:        host.Name,
			Address:     host.Address,
			Port:        host.Port,
			Enabled:     host.Enabled,
			GroupIDs:    append(host.GroupIDs, id),
			Environment: host.Environment,
		}); err != nil {
			return nil, nil, web.CommonErrors.ServerError
		}
	}
	for _, hostID := range removedHosts {
		host := HostCache.ByID(hostID)
		if !stringSliceContains(id, host.GroupIDs) {
			continue
		}

		if _, err := HostStore.EditHost(host, editHostParameters{
			Name:        host.Name,
			Address:     host.Address,
			Port:        host.Port,
			Enabled:     host.Enabled,
			GroupIDs:    filterStringSlice(id, host.GroupIDs),
			Environment: host.Environment,
		}); err != nil {
			return nil, nil, web.CommonErrors.ServerError
		}
	}

	hosts, err := group.Hosts()
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	return hosts, nil, nil
}

func (h *handle) GroupGetScripts(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	group := GroupCache.ByID(id)
	if group == nil {
		return nil, nil, web.ValidationError("No group with ID %s", id)
	}

	scripts, err := group.Scripts()
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}
	sort.Slice(scripts, func(i int, j int) bool {
		return scripts[i].Name < scripts[j].Name
	})

	return scripts, nil, nil
}

func (h *handle) GroupGetSchedules(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	schedules := ScheduleStore.AllSchedulesForGroup(id)
	sort.Slice(schedules, func(i int, j int) bool {
		return schedules[i].Name < schedules[j].Name
	})

	return schedules, nil, nil
}

func (h *handle) GroupNew(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	params := newGroupParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	group, err := GroupStore.NewGroup(params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.GroupAdded(group, session.Username)

	return group, nil, nil
}

func (h *handle) GroupEdit(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Parameters["id"]

	group := GroupCache.ByID(id)
	if group == nil {
		return nil, nil, web.ValidationError("No group with ID %s", id)
	}

	params := editGroupParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	group, err := GroupStore.EditGroup(group, params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.GroupModified(group, session.Username)

	return group, nil, nil
}

func (h *handle) GroupDelete(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Parameters["id"]

	group := GroupCache.ByID(id)
	if group == nil {
		return nil, nil, web.ValidationError("No group with ID %s", id)
	}

	if err := GroupStore.DeleteGroup(group); err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.GroupDeleted(group, session.Username)

	return true, nil, nil
}
