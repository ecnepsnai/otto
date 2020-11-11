package server

import (
	"regexp"

	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/security"
	"github.com/ecnepsnai/web"
)

func (h *handle) Register(request web.Request) (interface{}, *web.Error) {
	if !Options.Register.Enabled {
		log.Warn("Register is not enabled - rejecting request")
		return nil, web.CommonErrors.NotFound
	}

	r := otto.RegisterRequest{}
	if err := request.Decode(&r); err != nil {
		return nil, err
	}

	if Options.Register.PSK != r.PSK {
		log.Error("Invalid PSK for register request")
		return nil, web.CommonErrors.Unauthorized
	}

	existing, err := HostStore.HostWithAddress(r.Hostname)
	if err != nil {
		return nil, web.CommonErrors.ServerError
	}
	if existing != nil {
		return nil, web.ValidationError("Host with address '%s' already registered", r.Hostname)
	}

	groupID := Options.Register.DefaultGroupID
	for _, rule := range Options.Register.Rules {
		pattern, err := regexp.Compile(rule.Pattern)
		if err != nil {
			log.Error("Invalid registration rule regex: %s: %s", rule.Pattern, err.Error())
			continue
		}

		switch rule.Property {
		case RegisterRulePropertyUname:
			if pattern.MatchString(r.Uname) {
				log.Debug("Client uname '%s' matches pattern '%s'", r.Uname, rule.Property)
				groupID = rule.GroupID
				break
			}
		case RegisterRulePropertyHostname:
			if pattern.MatchString(r.Hostname) {
				log.Debug("Client hostname '%s' matches pattern '%s'", r.Hostname, rule.Property)
				groupID = rule.GroupID
				break
			}
		}
	}

	psk := security.RandomString(32)
	host, err := HostStore.NewHost(newHostParameters{
		Name:     r.Hostname,
		Address:  r.Address,
		Port:     r.Port,
		PSK:      psk,
		GroupIDs: []string{groupID},
	})
	if err != nil {
		log.Error("Error adding new host '%s': %s", r.Hostname, err.Message)
		return nil, web.CommonErrors.ServerError
	}
	log.Info("Registered new host '%s' -> '%s'", r.Hostname, host.ID)
	return otto.RegisterResponse{
		PSK: psk,
	}, nil
}
