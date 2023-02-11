package server

import (
	"strconv"

	"github.com/ecnepsnai/web"
)

func (h *handle) EventsGet(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	countStr := sliceFirst(request.HTTP.URL.Query()["c"])
	if countStr == "" {
		return nil, nil, web.ValidationError("Must specify max number of events")
	}

	count, cerr := strconv.Atoi(countStr)
	if cerr != nil {
		return nil, nil, web.ValidationError("Invalid count")
	}

	events, err := EventStore.LastEvents(count)
	if err != nil {
		if err.Server {
			return nil, nil, web.CommonErrors.ServerError
		}
		return nil, nil, web.ValidationError(err.Message)
	}

	return events, nil, nil
}
