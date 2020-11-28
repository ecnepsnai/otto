package server

import (
	"strconv"

	"github.com/ecnepsnai/web"
)

func (h *handle) EventsGet(request web.Request) (interface{}, *web.Error) {
	cStr := sliceFirst(request.HTTP.URL.Query()["c"])
	if cStr == "" {
		return nil, web.ValidationError("Must specify max number of events")
	}

	count, cerr := strconv.Atoi(cStr)
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
