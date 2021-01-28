package server

import (
	"github.com/ecnepsnai/web"
)

func (h *handle) ScheduleList(request web.Request) (interface{}, *web.Error) {
	return ScheduleStore.AllSchedules(), nil
}

func (h *handle) ScheduleGet(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedule := ScheduleStore.ScheduleWithID(id)
	if schedule == nil {
		return nil, web.ValidationError("No schedule with ID %s", id)
	}

	return schedule, nil
}

func (h *handle) ScheduleGetReports(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")
	return ScheduleReportStore.GetReportsForSchedule(id), nil
}

func (h *handle) ScheduleGetGroups(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedule := ScheduleStore.ScheduleWithID(id)
	if schedule == nil {
		return nil, web.ValidationError("No schedule with ID %s", id)
	}

	groups, err := schedule.Scope.Groups()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return groups, nil
}

func (h *handle) ScheduleGetHosts(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedule := ScheduleStore.ScheduleWithID(id)
	if schedule == nil {
		return nil, web.ValidationError("No schedule with ID %s", id)
	}

	hosts, err := schedule.Scope.Hosts()
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return hosts, nil
}

func (h *handle) ScheduleGetScript(request web.Request) (interface{}, *web.Error) {
	id := request.Params.ByName("id")

	schedule := ScheduleStore.ScheduleWithID(id)
	if schedule == nil {
		return nil, web.ValidationError("No schedule with ID %s", id)
	}

	script := ScriptStore.ScriptWithID(schedule.ScriptID)
	return script, nil
}

func (h *handle) ScheduleNew(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	params := newScheduleParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	schedule, err := ScheduleStore.NewSchedule(params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.ScheduleAdded(schedule, session.Username)

	return schedule, nil
}

func (h *handle) ScheduleEdit(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Params.ByName("id")

	schedule := ScheduleStore.ScheduleWithID(id)
	if schedule == nil {
		return nil, web.ValidationError("No schedule with ID %s", id)
	}

	params := editScheduleParameters{}
	if err := request.Decode(&params); err != nil {
		return nil, err
	}

	schedule, err := ScheduleStore.EditSchedule(schedule, params)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.ScheduleModified(schedule, session.Username)

	return schedule, nil
}

func (h *handle) ScheduleDelete(request web.Request) (interface{}, *web.Error) {
	session := request.UserData.(*Session)

	id := request.Params.ByName("id")

	schedule := ScheduleStore.ScheduleWithID(id)
	if schedule == nil {
		return nil, web.ValidationError("No schedule with ID %s", id)
	}

	if err := ScheduleStore.DeleteSchedule(schedule); err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	EventStore.ScheduleDeleted(schedule, session.Username)

	return true, nil
}
