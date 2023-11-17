package server

import (
	"os"
	"runtime"

	"github.com/ecnepsnai/web"
)

func (h *handle) State(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	type runtimeType struct {
		ServerFQDN    string
		Version       string
		BuildDate     string
		BuildRevision string
		Config        string
		Verbose       bool
	}
	type stateType struct {
		Runtime  runtimeType
		User     *User
		Warnings []string
		Options  *OttoOptions
	}

	hostname, _ := os.Hostname()
	user := request.UserData.(*Session).User()

	s := stateType{
		Runtime: runtimeType{
			ServerFQDN:    hostname,
			Version:       Version,
			BuildDate:     BuildDate,
			BuildRevision: BuildRevision,
			Config:        runtime.GOOS + "_" + runtime.GOARCH,
			Verbose:       verboseEnabled,
		},
		User:     user,
		Warnings: []string{},
	}
	if user.Permissions.CanModifySystem {
		s.Options = Options
	}

	if user.Username == defaultUser.Username {
		if ShadowStore.Compare(user.Username, []byte(defaultUser.Password)) {
			s.Warnings = append(s.Warnings, "default_user_password")
		}
	}

	return s, nil, nil
}

func (h *handle) Stats(request web.Request) (interface{}, *web.APIResponse, *web.Error) {
	return Stats.GetCounterValues(), nil, nil
}
