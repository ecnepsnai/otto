package server

import (
	"os"
	"runtime"

	"github.com/ecnepsnai/security"
	"github.com/ecnepsnai/web"
)

func (h *handle) State(request web.Request) (interface{}, *web.Error) {
	type runtimeType struct {
		ServerFQDN string
		Version    string
		Config     string
	}
	type stateType struct {
		Runtime  runtimeType
		User     *User
		Enums    map[string]interface{}
		Warnings []string
		Options  *OttoOptions
	}

	hostname, _ := os.Hostname()
	user := request.UserData.(*Session).User()

	s := stateType{
		Runtime: runtimeType{
			ServerFQDN: hostname,
			Version:    ServerVersion,
			Config:     runtime.GOOS + "_" + runtime.GOARCH,
		},
		User:     user,
		Enums:    AllEnums,
		Warnings: []string{},
		Options:  Options,
	}

	if user.Username == defaultUser.Username {
		delay := security.FailDelay
		security.FailDelay = 0
		if user.PasswordHash.Compare([]byte(defaultUser.Password)) {
			s.Warnings = append(s.Warnings, "default_user_password")
		}
		security.FailDelay = delay
	}

	return s, nil
}
