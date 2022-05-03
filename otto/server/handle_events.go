package server

import (
	"strconv"

	"github.com/ecnepsnai/web"
)

func (h *handle) EventsGet(request web.Request) (interface{}, *web.Error) {
	countStr := sliceFirst(request.HTTP.URL.Query()["c"])
	if countStr == "" {
		return nil, web.ValidationError("Must specify max number of events")
	}

	count, cerr := strconv.Atoi(countStr)
	if cerr != nil {
		return nil, web.ValidationError("Invalid count")
	}

	events, err := EventStore.LastEvents(count)
	if err != nil {
		if err.Server {
			return nil, web.CommonErrors.ServerError
		}
		return nil, web.ValidationError(err.Message)
	}

	return events, nil
}
