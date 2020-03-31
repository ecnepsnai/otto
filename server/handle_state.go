package server

import (
	"os"

	"github.com/ecnepsnai/web"
)

func (h *handle) State(request web.Request) (interface{}, *web.Error) {
	type state struct {
		User       *User
		ServerFQDN string
		Version    string
		Enums      map[string]interface{}
	}

	hostname, _ := os.Hostname()
	user := request.UserData.(*Session).User()

	return state{
		User:       user,
		ServerFQDN: hostname,
		Version:    ServerVersion,
		Enums:      AllEnums,
	}, nil
}
