package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ecnepsnai/otto/shared/otto"
	"github.com/ecnepsnai/secutil"
	"github.com/ecnepsnai/web"
)

func (v *view) Register(request web.Request, writer web.Writer) web.HTTPResponse {
	if !Options.Register.Enabled {
		log.PWarn("Rejected registration request", map[string]interface{}{
			"remote_addr": request.HTTP.RemoteAddr,
			"reason":      "registration disabled",
		})
		return web.HTTPResponse{
			Status: 404,
		}
	}

	if requestProtocolVersion := request.HTTP.Header.Get("X-OTTO-PROTO-VERSION"); fmt.Sprintf("%d", otto.ProtocolVersion) != requestProtocolVersion {
		log.PWarn("Rejected registration request", map[string]interface{}{
			"remote_addr":            request.HTTP.RemoteAddr,
			"reason":                 "unsupported otto protocol version",
			"agent_protocol_version": requestProtocolVersion,
		})
		return web.HTTPResponse{
			Status: 400,
		}
	}

	encryptedData, erro := io.ReadAll(request.HTTP.Body)
	if erro != nil {
		return web.HTTPResponse{
			Status: 400,
		}
	}
	decryptedData, erro := secutil.Encryption.AES_256_GCM.Decrypt(encryptedData, Options.Register.Key)
	if erro != nil {
		EventStore.HostRegisterIncorrectKey(request.HTTP.RemoteAddr)
		log.PWarn("Rejected registration request", map[string]interface{}{
			"remote_addr": request.HTTP.RemoteAddr,
			"reason":      "incorrect register key",
		})
		return web.HTTPResponse{
			Status: 403,
		}
	}

	r := otto.RegisterRequest{}
	if err := json.Unmarshal(decryptedData, &r); err != nil {
		return web.HTTPResponse{
			Status: 400,
		}
	}
	hostIP := request.ClientIPAddress().String()
	log.PDebug("Incoming host registration request", map[string]interface{}{
		"HostIP":              hostIP,
		"AgentIdentity":       r.AgentIdentity,
		"Port":                r.Port,
		"Hostname":            r.Properties.Hostname,
		"KernelName":          r.Properties.KernelName,
		"KernelVersion":       r.Properties.KernelVersion,
		"DistributionName":    r.Properties.DistributionName,
		"DistributionVersion": r.Properties.DistributionVersion,
	})

	if existing := HostStore.HostWithName(r.Properties.Hostname); existing != nil {
		log.PWarn("Rejected registration request", map[string]interface{}{
			"reason":   "duplicate hostname",
			"hostname": r.Properties.Hostname,
		})
		return web.HTTPResponse{
			Status: 400,
		}
	}

	if existing := HostStore.HostWithAddress(hostIP); existing != nil {
		log.PWarn("Rejected registration request", map[string]interface{}{
			"reason":  "duplicate address",
			"address": hostIP,
		})
		return web.HTTPResponse{
			Status: 400,
		}
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
	host, err := HostStore.NewHost(newHostParameters{
		Name:          r.Properties.Hostname,
		Address:       hostIP,
		Port:          r.Port,
		AgentIdentity: r.AgentIdentity,
		GroupIDs:      []string{groupID},
	})
	if err != nil {
		log.PError("Error adding new host", map[string]interface{}{
			"hostname": r.Properties.Hostname,
			"address":  hostIP,
			"error":    err.Message,
		})
		return web.HTTPResponse{
			Status: 500,
		}
	}
	EventStore.HostRegisterSuccess(host, r, matchedRule)

	serverId, idErr := IdentityStore.Get(host.ID)
	if idErr != nil {
		log.PError("Error getting identity for host", map[string]interface{}{
			"host_id": host.ID,
			"error":   idErr.Error(),
		})
		return web.HTTPResponse{
			Status: 500,
		}
	}
	if serverId == nil {
		log.PError("No server identity for host", map[string]interface{}{"host_id": host.ID})
		return web.HTTPResponse{
			Status: 500,
		}
	}

	responseData, erro := json.Marshal(otto.RegisterResponse{
		ServerIdentity: serverId.PublicKeyString(),
	})
	if erro != nil {
		return web.HTTPResponse{
			Status: 500,
		}
	}
	encryptedResponse, erro := secutil.Encryption.AES_256_GCM.Encrypt(responseData, Options.Register.Key)
	if erro != nil {
		return web.HTTPResponse{
			Status: 500,
		}
	}
	return web.HTTPResponse{
		Status:      200,
		Reader:      io.NopCloser(bytes.NewReader(encryptedResponse)),
		ContentType: "application/octet-stream",
		Headers: map[string]string{
			"Content-Length":       fmt.Sprintf("%d", len(encryptedResponse)),
			"X-OTTO-PROTO-VERSION": fmt.Sprintf("%d", otto.ProtocolVersion),
		},
	}
}
