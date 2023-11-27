package server

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto/server/environ"
	"github.com/ecnepsnai/otto/shared/otto"
	"github.com/ecnepsnai/secutil"
)

// ScriptResult describes a script result
type ScriptResult struct {
	ScriptID    string
	Duration    time.Duration
	Environment []environ.Variable
	Result      otto.ScriptResult
	RunError    string
}

type hostConnection struct {
	Host     *Host
	Address  string
	Identity *otto.Identity
	Conn     *otto.Connection
}

// connect will open a connection to the Otto agent on the host
func (host *Host) connect() (*hostConnection, error) {
	address := fmt.Sprintf("%s:%d", host.Address, host.Port)
	log.Debug("Connecting to host %s", address)

	timeout := time.Duration(Options.Network.Timeout) * time.Second
	network := "tcp"
	if Options.Network.ForceIPVersion == IPVersionOptionIPv4 {
		network = "tcp4"
	} else if Options.Network.ForceIPVersion == IPVersionOptionIPv6 {
		network = "tcp6"
	}

	id, err := IdentityStore.Get(host.ID)
	if err != nil {
		log.PError("Error getting identity for host", map[string]interface{}{
			"host_id": host.ID,
			"error":   err.Error(),
		})
		return nil, fmt.Errorf("identity: %s", err.Error())
	}
	if id == nil {
		log.PError("No server identity for host", map[string]interface{}{"host_id": host.ID})
		return nil, fmt.Errorf("no identity")
	}

	connection, err := otto.Dial(otto.DialOptions{
		Network:          network,
		Address:          address,
		Identity:         id.Signer(),
		TrustedPublicKey: host.Trust.TrustedIdentity,
		Timeout:          timeout,
	})
	log.Debug("dialed host %s", host.ID)
	if err != nil {
		if strings.Contains(err.Error(), "unknown public key:") {
			parts := strings.Split(err.Error(), " ")
			keyHex := parts[len(parts)-1]
			key, herr := hex.DecodeString(keyHex)
			if herr != nil {
				return nil, err
			}
			pendingKey := base64.StdEncoding.EncodeToString(key)
			host.Trust.UntrustedIdentity = pendingKey
			HostStore.UpdateHostTrust(host.ID, host.Trust)
			log.PInfo("Recorded new identity from agent", map[string]interface{}{
				"agent":    host.ID,
				"identity": pendingKey,
			})
		}

		heartbeatStore.UpdateHostReachability(host, false)
		log.Error("Error connecting to host '%s': %s", address, err.Error())
		return nil, err
	}
	return &hostConnection{
		Host:     host,
		Address:  fmt.Sprintf("%s:%d", host.Address, host.Port),
		Identity: id,
		Conn:     connection,
	}, nil
}

// Close will close the connection to the Otto host
func (hc *hostConnection) Close() {
	hc.Conn.Close()
}

func (conn *hostConnection) UploadFile(attachment Attachment) error {
	fileInfo, err := attachment.FileInfo()
	if err != nil {
		log.PError("Error uploading file", map[string]interface{}{
			"attachment_id": attachment.ID,
			"error":         err.Error(),
		})
		return err
	}
	reader, err := attachment.Reader()
	if err != nil {
		log.PError("Error uploading file", map[string]interface{}{
			"attachment_id": attachment.ID,
			"error":         err.Error(),
		})
		return err
	}

	if err := conn.Conn.TriggerActionUploadFile(*fileInfo, reader); err != nil {
		log.PError("Error uploading file", map[string]interface{}{
			"attachment_id": attachment.ID,
			"error":         err.Error(),
		})
		return err
	}

	return nil
}

// TODO: change this to read the script as a file
func (conn *hostConnection) RunScript(scriptInfo otto.ScriptInfo, scriptData []byte, actionOutput func(stdout, stderr []byte), cancel chan bool) (*otto.MessageActionResult, error) {
	result, err := conn.Conn.TriggerActionRunScript(scriptInfo, io.NopCloser(bytes.NewReader(scriptData)), actionOutput, cancel)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Ping ping the host
func (host *Host) Ping() error {
	conn, err := host.connect()
	if err != nil {
		log.Error("Error sending heartbeat request to host '%s': %s", host.ID, err.Error())
		heartbeatStore.UpdateHostReachability(host, false)
		return err
	}
	defer conn.Close()

	nonce := secutil.RandomString(8)
	reply, err := conn.Conn.SendHeartbeat(otto.MessageHeartbeatRequest{Version: Version, Nonce: nonce})
	if err != nil {
		log.PError("Error sending heartbeat request to host", map[string]interface{}{
			"host_id": host.ID,
			"error":   err.Error(),
		})
		heartbeatStore.UpdateHostReachability(host, false)
		return err
	}
	if reply.Nonce != nonce {
		log.PError("Unexpected nonce in heartbeat reply", map[string]interface{}{
			"host_id":        host.ID,
			"expected_nonce": nonce,
			"actual_nonce":   reply.Nonce,
		})
		heartbeatStore.UpdateHostReachability(host, false)
		return fmt.Errorf("invalid nonce")
	}
	heartbeatStore.RegisterHeartbeatReply(host, *reply)

	return nil
}

// RunScript run the script on the host. Error will only ever be populated with internal server
// errors, such as being unable to read from the database.
func (host *Host) RunScript(script *Script, scriptOutput func(stdout, stderr []byte), cancel chan bool) (*ScriptResult, error) {
	start := time.Now()
	log.PInfo("Running script on host", map[string]interface{}{
		"host_id":   host.ID,
		"script_id": script.ID,
	})

	scriptRequest := script.ScriptInfo()

	variables := host.environmentVariablesForScript(script)
	scriptRequest.Environment = environ.Map(variables)
	log.Debug("Environ: %s", scriptRequest.Environment)

	attachments, aerr := script.Attachments()
	if aerr != nil {
		return nil, aerr.Error
	}

	conn, err := host.connect()
	if err != nil {
		log.PError("Error running script on host", map[string]interface{}{
			"script_id": script.ID,
			"host_id":   host.ID,
			"error":     err.Error(),
		})
		heartbeatStore.UpdateHostReachability(host, false)
		return nil, err
	}

	// Pre-execution files
	for _, attachment := range attachments {
		if attachment.AfterScript {
			continue
		}

		log.PInfo("Uploading script attachment", map[string]interface{}{
			"script_id":     script.ID,
			"attachment_id": attachment.ID,
			"host_id":       host.ID,
		})
		if err := conn.UploadFile(attachment); err != nil {
			log.PError("Error running script on host", map[string]interface{}{
				"script_id": script.ID,
				"host_id":   host.ID,
				"error":     err.Error(),
			})
			return nil, err
		}
	}

	// Execute script
	log.PInfo("Executing script", map[string]interface{}{
		"script_id": script.ID,
		"host_id":   host.ID,
	})
	result, err := conn.RunScript(scriptRequest, []byte(script.Script), scriptOutput, cancel)
	if result == nil && err == nil {
		err = fmt.Errorf("unexpected end of connection")
	}
	if err != nil {
		log.PError("Error running script on host", map[string]interface{}{
			"host_id":   host.ID,
			"script_id": script.ID,
			"error":     err.Error(),
		})
		return &ScriptResult{
			ScriptID:    script.ID,
			Duration:    time.Since(start),
			Environment: variables,
			Result: otto.ScriptResult{
				Success: false,
			},
			RunError: err.Error(),
		}, nil
	}

	if !result.ScriptResult.Success {
		log.PError("Error running script on host", map[string]interface{}{
			"host_id":   host.ID,
			"script_id": script.ID,
			"error":     result.ScriptResult.ExecError,
		})
		return &ScriptResult{
			ScriptID:    script.ID,
			Duration:    time.Since(start),
			Environment: variables,
			Result:      result.ScriptResult,
		}, nil
	}

	// Post-execution files
	for _, attachment := range attachments {
		if !attachment.AfterScript {
			continue
		}

		log.PInfo("Uploading script attachment", map[string]interface{}{
			"script_id":     script.ID,
			"attachment_id": attachment.ID,
			"host_id":       host.ID,
		})
		if err := conn.UploadFile(attachment); err != nil {
			log.PError("Error running script on host", map[string]interface{}{
				"script_id": script.ID,
				"host_id":   host.ID,
				"error":     err.Error(),
			})
			return nil, err
		}
	}

	heartbeatStore.UpdateHostReachability(host, true)

	// After execution actions
	switch script.AfterExecution {
	case AgentActionReloadConfig:
		err = conn.Conn.TriggerActionReloadConfig()
	case AgentActionExitAgent:
		err = conn.Conn.TriggerActionExitAgent()
	case AgentActionReboot:
		err = conn.Conn.TriggerActionReboot()
	case AgentActionShutdown:
		err = conn.Conn.TriggerActionShutdown()
	case "":
		// Noop
		err = nil
	default:
		log.PError("Unknown post-execution action", map[string]interface{}{
			"action":    script.AfterExecution,
			"script_id": script.ID,
		})
		return &ScriptResult{
			ScriptID:    script.ID,
			Duration:    time.Since(start),
			Environment: variables,
			Result: otto.ScriptResult{
				Success: false,
			},
			RunError: fmt.Sprintf("unknown post-execution action %s", script.AfterExecution),
		}, nil
	}
	if err != nil {
		log.PError("Error running script post-execution action on host", map[string]interface{}{
			"host_id":   host.ID,
			"script_id": script.ID,
			"error":     result.ScriptResult.ExecError,
		})
		return &ScriptResult{
			ScriptID:    script.ID,
			Duration:    time.Since(start),
			Environment: variables,
			Result: otto.ScriptResult{
				Success: false,
			},
			RunError: err.Error(),
		}, nil
	}

	log.PInfo("Finished running script on host", map[string]interface{}{
		"host_id":   host.ID,
		"script_id": script.ID,
		"elapsed":   time.Since(start).String(),
	})
	return &ScriptResult{
		ScriptID:    script.ID,
		Duration:    time.Since(start),
		Environment: variables,
		Result:      result.ScriptResult,
	}, nil
}

func (host *Host) environmentVariablesForScript(script *Script) []environ.Variable {
	variables := environ.Merge(staticEnvironment(), []environ.Variable{
		environ.New("OTTO_HOST_ADDRESS", host.Address),
		environ.New("OTTO_HOST_PORT", fmt.Sprintf("%d", host.Port)),
	})

	// 1. Global environment variables
	variables = environ.Merge(variables, Options.General.GlobalEnvironment)

	// 2. Script environment variables
	variables = environ.Merge(variables, script.Environment)

	// 3. Group environment variables
	groups, err := host.Groups()
	if err != nil {
		groups = []Group{}
	}
	for _, group := range groups {
		variables = environ.Merge(variables, group.Environment)
	}

	// 4. Host environment variables
	variables = environ.Merge(variables, host.Environment)

	if logtic.Log.Level == logtic.LevelDebug {
		varStr := make([]string, len(variables))
		for i, variable := range variables {
			if variable.Secret {
				varStr[i] = variable.Key + "='*****'"
			} else {
				varStr[i] = fmt.Sprintf("%s='%s'", variable.Key, variable.Value)
			}
		}
		log.Debug("Script variables: %s", strings.Join(varStr, " "))
	}

	return variables
}

// RotateIdentity will rotate the identity for both the server and client. Returns the server public key, client public
// key, or an error
func (host *Host) RotateIdentity() (string, string, error) {
	serverId, iderr := otto.NewIdentity()
	if iderr != nil {
		log.PError("Error generating new identity", map[string]interface{}{
			"error": iderr.Error(),
		})
		return "", "", iderr
	}

	conn, err := host.connect()
	if err != nil {
		heartbeatStore.UpdateHostReachability(host, false)
		log.Error("Error triggering action on host '%s': %s", host.ID, err.Error())
		return "", "", err
	}
	defer conn.Close()

	reply, err := conn.Conn.RotateIdentity(otto.MessageRotateIdentityRequest{
		PublicKey: serverId.PublicKeyString(),
	})
	if err != nil {
		log.PError("Error requesting agent update identity", map[string]interface{}{
			"host":  host.ID,
			"error": err.Error,
		})
		return "", "", err
	}
	if reply.Error != "" {
		log.PError("Error requesting agent update identity", map[string]interface{}{
			"host":  host.ID,
			"error": reply.Error,
		})
		return "", "", fmt.Errorf(reply.Error)
	}
	if reply.PublicKey == "" {
		log.PError("Error requesting agent update identity", map[string]interface{}{
			"host":  host.ID,
			"error": "no public key in response",
		})
		return "", "", fmt.Errorf("no public key in response")
	}

	agentPublicKey := reply.PublicKey

	IdentityStore.Set(host.ID, serverId)
	trust := HostTrust{
		TrustedIdentity: agentPublicKey,
		LastTrustUpdate: time.Now(),
	}
	if err := HostStore.UpdateHostTrust(host.ID, trust); err != nil {
		log.PError("Error updating host trust", map[string]interface{}{
			"host":  host.ID,
			"error": err.Message,
		})
		return "", "", err.Error
	}

	log.PInfo("Rotated host identities", map[string]interface{}{
		"host_id":    host.ID,
		"host_name":  host.Name,
		"server_pub": serverId.PublicKeyString(),
		"agent_pub":  agentPublicKey,
	})
	return serverId.PublicKeyString(), agentPublicKey, nil
}
