package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ecnepsnai/otto"
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
			"remote_addr":             request.HTTP.RemoteAddr,
			"reason":                  "unsupported otto protocol version",
			"client_protocol_version": requestProtocolVersion,
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

	existing := HostStore.HostWithAddress(r.Properties.Hostname)
	if existing != nil {
		log.PWarn("Rejected registration request", map[string]interface{}{
			"remote_addr": request.HTTP.RemoteAddr,
			"reason":      "duplicate hostname",
			"hostname":    r.Properties.Hostname,
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
	address := stripPortFromRemoteAddr(request.HTTP.RemoteAddr)
	host, err := HostStore.NewHost(newHostParameters{
		Name:           r.Properties.Hostname,
		Address:        address,
		Port:           r.Port,
		ClientIdentity: r.ClientIdentity,
		GroupIDs:       []string{groupID},
	})
	if err != nil {
		log.PError("Error adding new host", map[string]interface{}{
			"hostname": r.Properties.Hostname,
			"address":  address,
			"error":    err.Message,
		})
		return web.HTTPResponse{
			Status: 500,
		}
	}
	EventStore.HostRegisterSuccess(host, r, matchedRule)

	serverId := IdentityStore.Get(host.ID)
	if serverId == nil {
		log.PError("No server identity for host", map[string]interface{}{"host_id": host.ID})
		return web.HTTPResponse{
			Status: 500,
		}
	}

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

	responseData, erro := json.Marshal(otto.RegisterResponse{
		ServerIdentity: serverId.PublicKeyString(),
		Scripts:        scripts,
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
