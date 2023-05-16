package server

import (
	"fmt"
	"sort"

	"github.com/ecnepsnai/web"
)

func (h *handle) ScheduleList(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	schedules := ScheduleCache.All()
	sort.Slice(schedules, func(i int, j int) bool {
		return schedules[i].Name < schedules[j].Name
	})

	return schedules, nil, nil
}

func (h *handle) ScheduleGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	schedule := ScheduleCache.ByID(id)
	if schedule == nil {
		return nil, nil, web.ValidationError("No schedule with ID %s", id)
	}

	return schedule, nil, nil
}

func (h *handle) ScheduleGetReports(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]
	return ScheduleReportStore.GetReportsForSchedule(id), nil, nil
}

func (h *handle) ScheduleGetGroups(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	schedule := ScheduleCache.ByID(id)
	if schedule == nil {
		return nil, nil, web.ValidationError("No schedule with ID %s", id)
	}

	groups, err := schedule.Scope.Groups()
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}
	sort.Slice(groups, func(i int, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups, nil, nil
}

func (h *handle) ScheduleGetHosts(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	schedule := ScheduleCache.ByID(id)
	if schedule == nil {
		return nil, nil, web.ValidationError("No schedule with ID %s", id)
	}

	hosts, err := schedule.Scope.Hosts()
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

func (h *handle) ScheduleGetScript(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	id := request.Parameters["id"]

	schedule := ScheduleCache.ByID(id)
	if schedule == nil {
		return nil, nil, web.ValidationError("No schedule with ID %s", id)
	}

	script := ScriptStore.ScriptWithID(schedule.ScriptID)
	return script, nil, nil
}

func (h *handle) ScheduleNew(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)

	if !session.User().Permissions.CanModifySchedules {
		EventStore.UserPermissionDenied(session.User().Username, "Create new schedule")
		return nil, nil, web.ValidationError("Permission denied")
	}

	params := newScheduleParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	schedule, err := ScheduleStore.NewSchedule(params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.ScheduleAdded(schedule, session.Username)

	return schedule, nil, nil
}

func (h *handle) ScheduleEdit(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifySchedules {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Modify schedule %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	schedule := ScheduleCache.ByID(id)
	if schedule == nil {
		return nil, nil, web.ValidationError("No schedule with ID %s", id)
	}

	params := editScheduleParameters{}
	if err := request.DecodeJSON(&params); err != nil {
		return nil, nil, err
	}

	schedule, err := ScheduleStore.EditSchedule(schedule, params)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.ScheduleModified(schedule, session.Username)

	return schedule, nil, nil
}

func (h *handle) ScheduleDelete(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	session := request.UserData.(*Session)
	id := request.Parameters["id"]

	if !session.User().Permissions.CanModifySchedules {
		EventStore.UserPermissionDenied(session.User().Username, fmt.Sprintf("Delete schedule %s", id))
		return nil, nil, web.ValidationError("Permission denied")
	}

	schedule := ScheduleCache.ByID(id)
	if schedule == nil {
		return nil, nil, web.ValidationError("No schedule with ID %s", id)
	}

	if err := ScheduleStore.DeleteSchedule(schedule); err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	EventStore.ScheduleDeleted(schedule, session.Username)

	return true, nil, nil
}
