package server

import (
	"github.com/ecnepsnai/web"
)

func (h *handle) HeartbeatLast(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	return heartbeatStore.AllHeartbeats(), nil, nil
}
