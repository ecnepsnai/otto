package server

import (
	"sort"

	"github.com/ecnepsnai/web"
)

func (h *handle) GroupList(request web.Request) (interface{}, *web.Error) {
	groups := GroupStore.AllGroups()
	sort.Slice(groups, func(i int, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups, nil
}

func (h *handle) GroupGetMembership(request web.Request) (interface{}, *web.Error) {
	return GroupCache.Membership(), nil
}

func (h *handle) GroupGet(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	group := GroupStore.GroupWithID(id)
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	return group, nil
}

func (h *handle) GroupGetHosts(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	group := GroupStore.GroupWithID(id)
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
	sort.Slice(hosts, func(i int, j int) bool {
		return hosts[i].Name < hosts[j].Name
	})

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

	group := GroupStore.GroupWithID(id)
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
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
		host := HostStore.HostWithID(hostID)
		if stringSliceContains(id, host.GroupIDs) {
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
	for _, hostID := range removedHosts {
		host := HostStore.HostWithID(hostID)
		if !stringSliceContains(id, host.GroupIDs) {
			continue
		}

		if _, err := HostStore.EditHost(host, editHostParameters{
			Name:        host.Name,
			Address:     host.Address,
			Port:        host.Port,
			PSK:         host.PSK,
			Enabled:     host.Enabled,
			GroupIDs:    filterStringSlice(id, host.GroupIDs),
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

	group := GroupStore.GroupWithID(id)
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
	sort.Slice(scripts, func(i int, j int) bool {
		return scripts[i].Name < scripts[j].Name
	})

	return scripts, nil
}

func (h *handle) GroupGetSchedules(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedules := ScheduleStore.AllSchedulesForGroup(id)
	sort.Slice(schedules, func(i int, j int) bool {
		return schedules[i].Name < schedules[j].Name
	})

	return schedules, nil
}

func (h *handle) GroupNew(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

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

	EventStore.GroupAdded(group, session.Username)

	return group, nil
}

func (h *handle) GroupEdit(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Params.ByName("id")

	group := GroupStore.GroupWithID(id)
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	params := editGroupParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	group, err := GroupStore.EditGroup(group, params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.GroupModified(group, session.Username)

	return group, nil
}

func (h *handle) GroupDelete(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Params.ByName("id")

	group := GroupStore.GroupWithID(id)
	if group == nil {
		return nil, web.ValidationError("No group with ID %s", id)
	}

	if err := GroupStore.DeleteGroup(group); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.GroupDeleted(group, session.Username)

	return true, nil
}
