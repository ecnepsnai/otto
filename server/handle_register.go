package server

import (
	"fmt"
	"time"

	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/web"
)

func (h *handle) Register(request web.Request) (interface{}, *web.Error) {
	if !Options.Register.Enabled {
		log.PWarn("Rejected registration request", map[string]interface{}{
			"remote_addr": request.HTTP.RemoteAddr,
			"reason":      "registration disabled",
		})
		return nil, web.CommonErrors.NotFound
	}

	if requestProtocolVersion := request.HTTP.Header.Get("X-OTTO-PROTO-VERSION"); fmt.Sprintf("%d", otto.ProtocolVersion) != requestProtocolVersion {
		log.PWarn("Rejected registration request", map[string]interface{}{
			"remote_addr":             request.HTTP.RemoteAddr,
			"reason":                  "unsupported otto protocol version",
			"client_protocol_version": requestProtocolVersion,
		})
		return nil, web.ValidationError("Unsupported protocol version %s", requestProtocolVersion)
	}

	r := otto.RegisterRequest{}
	if err := request.Decode(&r); err != nil {
		return nil, err
	}

	if Options.Register.Key != r.Key {
		EventStore.HostRegisterIncorrectKey(request.HTTP.RemoteAddr, r)
		log.PWarn("Rejected registration request", map[string]interface{}{
			"remote_addr": request.HTTP.RemoteAddr,
			"reason":      "incorrect register key",
		})
		return nil, web.CommonErrors.Unauthorized
	}

	existing := HostStore.HostWithAddress(r.Properties.Hostname)
	if existing != nil {
		log.PWarn("Rejected registration request", map[string]interface{}{
			"remote_addr": request.HTTP.RemoteAddr,
			"reason":      "duplicate hostname",
			"hostname":    r.Properties.Hostname,
		})
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

	// Use the address that sent this request as the address for the new host, but first strip the port
	// RemoteAddr will always have a port, but may be a wrapped IPv6 address
	address := stripPortFromRemoteAddr(request.HTTP.RemoteAddr)

	psk := newHostPSK()
	host, err := HostStore.NewHost(newHostParameters{
		Name:          r.Properties.Hostname,
		Address:       address,
		Port:          r.Port,
		PSK:           psk,
		LastPSKRotate: time.Now(),
		GroupIDs:      []string{groupID},
	})
	if err != nil {
		log.PError("Error adding new host", map[string]interface{}{
			"hostname": r.Properties.Hostname,
			"address":  address,
			"error":    err.Message,
		})
		return nil, web.CommonErrors.ServerError
	}
	log.PInfo("Registered new host", map[string]interface{}{
		"name":    host.Name,
		"address": host.Address,
		"port":    host.Port,
	})
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
