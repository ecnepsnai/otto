package server

import (
	"fmt"
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

	existing, err := HostStore.HostWithAddress(r.Properties.Hostname)
	if err != nil {
		return nil, web.CommonErrors.ServerError
	}
	if existing != nil {
		return nil, web.ValidationError("Host with address '%s' already registered", r.Properties.Hostname)
	}

	groupID := Options.Register.DefaultGroupID
	var matchedRule *RegisterRule
	for _, rule := range RegisterRuleStore.AllRules() {
		pattern, err := regexp.Compile(rule.Pattern)
		if err != nil {
			log.Error("Invalid registration rule regex: %s: %s", rule.Pattern, err.Error())
			continue
		}

		switch rule.Property {
		case RegisterRulePropertyHostname:
			if pattern.MatchString(r.Properties.Hostname) {
				log.Debug("Property %s matches client value '%s'", RegisterRulePropertyHostname, r.Properties.Hostname)
				groupID = rule.GroupID
				matchedRule = &rule
			}
			break
		case RegisterRulePropertyKernelName:
			if pattern.MatchString(r.Properties.KernelName) {
				log.Debug("Property %s matches client value '%s'", RegisterRulePropertyKernelName, r.Properties.KernelName)
				groupID = rule.GroupID
				matchedRule = &rule
			}
			break
		case RegisterRulePropertyKernelVersion:
			if pattern.MatchString(r.Properties.KernelVersion) {
				log.Debug("Property %s matches client value '%s'", RegisterRulePropertyKernelVersion, r.Properties.KernelVersion)
				groupID = rule.GroupID
				matchedRule = &rule
			}
			break
		case RegisterRulePropertyDistributionName:
			if pattern.MatchString(r.Properties.DistributionName) {
				log.Debug("Property %s matches client value '%s'", RegisterRulePropertyDistributionName, r.Properties.DistributionName)
				groupID = rule.GroupID
				matchedRule = &rule
			}
			break
		case RegisterRulePropertyDistributionVersion:
			if pattern.MatchString(r.Properties.DistributionVersion) {
				log.Debug("Property %s matches client value '%s'", RegisterRulePropertyDistributionVersion, r.Properties.DistributionVersion)
				groupID = rule.GroupID
				matchedRule = &rule
			}
			break
		}
	}

	psk := security.RandomString(32)
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
