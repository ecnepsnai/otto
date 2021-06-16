package server

import (
	"fmt"
	"time"

	"github.com/ecnepsnai/otto"
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

	if Options.Register.Key != r.Key {
		EventStore.HostRegisterIncorrectKey(r)
		log.Error("Invalid Key for register request")
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

	psk := newHostPSK()
	host, err := HostStore.NewHost(newHostParameters{
		Name:          r.Properties.Hostname,
		Address:       r.Address,
		Port:          r.Port,
		PSK:           psk,
		LastPSKRotate: time.Now(),
		GroupIDs:      []string{groupID},
	})
	if err != nil {
		log.Error("Error adding new host '%s': %s", r.Properties.Hostname, err.Message)
		return nil, web.CommonErrors.ServerError
	}
	log.Info("Registered new host '%s' -> '%s'", r.Properties.Hostname, host.ID)
	EventStore.HostRegisterSuccess(host, r, matchedRule)

	scripts := []otto.Script{}
	if Options.Register.RunScriptsOnRegister {
		for _, scriptID := range host.Scripts() {
			script := ScriptStore.ScriptWithID(scriptID.ScriptID)
			scriptRequest, err := script.OttoScript()
			if err != nil {
				continue
			}
			scripts = append(scripts, *scriptRequest)
		}
	}

	return otto.RegisterResponse{
		PSK:     psk,
		Scripts: scripts,
	}, nil
}
