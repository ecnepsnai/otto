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
		Warnings   []string
	}

	hostname, _ := os.Hostname()
	user := request.UserData.(*Session).User()

	s := state{
		User:       user,
		ServerFQDN: hostname,
		Version:    ServerVersion,
		Enums:      AllEnums,
		Warnings:   []string{},
	}

	if user.Username == defaultUser.Username {
		if user.PasswordHash.Compare(defaultUser.Password) {
			s.Warnings = append(s.Warnings, "default_user_password")
		}
	}

	return s, nil
}
