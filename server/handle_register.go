package server

import (
	"fmt"

	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/secutil"
	"github.com/ecnepsnai/web"
)

func (h *handle) Register(request web.Request) (interface{}, *web.Error) {
	if !Options.Register.Enabled {
		log.Warn("Register is not enabled - rejecting request")
		return nil, web.CommonErrors.NotFound
	}

	if requestProtocolVersion := request.HTTP.Header.Get("X-OTTO-PROTO-VERSION"); fmt.Sprintf("%d", otto.ProtocolVersion) != requestProtocolVersion {
		return nil, web.ValidationError("Unsupported protocol version %s", requestProtocolVersion)
	}

	r := otto.RegisterRequest{}
	if err := request.Decode(&r); err != nil {
		return nil, err
	}

	if Options.Register.PSK != r.PSK {
		EventStore.HostRegisterIncorrectPSK(r)
		log.Error("Invalid PSK for register request")
		return nil, web.CommonErrors.Unauthorized
	}

	existing := HostStore.HostWithAddress(r.Properties.Hostname)
	if existing != nil {
		return nil, web.ValidationError("Host with address '%s' already registered", r.Properties.Hostname)
	}

	groupID := Options.Register.DefaultGroupID
	var matchedRule *RegisterRule
	for _, rule := range RegisterRuleStore.AllRules() {
		if !rule.Matches(r.Properties) {
			continue
		}

		groupID = rule.GroupID
		matchedRule = &rule
		break
	}

	psk := secutil.RandomString(32)
	host, err := HostStore.NewHost(newHostParameters{
		Name:     r.Properties.Hostname,
		Address:  r.Address,
		Port:     r.Port,
		PSK:      psk,
		GroupIDs: []string{groupID},
	})
	if err != nil {
		log.Error("Error adding new host '%s': %s", r.Properties.Hostname, err.Message)
		return nil, web.CommonErrors.ServerError
	}
	log.Info("Registered new host '%s' -> '%s'", r.Properties.Hostname, host.ID)
	EventStore.HostRegisterSuccess(host, r, matchedRule)
	return otto.RegisterResponse{
		PSK: psk,
	}, nil
}
