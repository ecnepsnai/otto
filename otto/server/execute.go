package server

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/otto/server/environ"
)

// ScriptResult describes a script result
type ScriptResult struct {
	ScriptID    string
	Duration    time.Duration
	Environment []environ.Variable
	Result      otto.ScriptResult
	RunError    string
}

var clientActionMap = map[string]uint32{
	ClientActionRunScript:      otto.ActionRunScript,
	ClientActionExitClient:     otto.ActionUploadFileAndExitClient,
	ClientActionReboot:         otto.ActionReboot,
	ClientActionShutdown:       otto.ActionShutdown,
	ClientActionUpdateIdentity: otto.ActionUpdateIdentity,
}

type hostConnection struct {
	Host     *Host
	Address  string
	Identity otto.Identity
	Conn     *otto.Connection
}

// connect will open a connection to the Otto client on the host
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

	id := IdentityStore.Get(host.ID)
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
			pendingKey := base64.RawStdEncoding.EncodeToString(key)
			host.Trust.UntrustedIdentity = pendingKey
			HostStore.UpdateHostTrust(host.ID, host.Trust)
			log.PInfo("Recorded new identity from client", map[string]interface{}{
				"client":   host.ID,
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

// TriggerAction will trigger the given action on the host
func (host *Host) TriggerAction(action otto.MessageTriggerAction, actionOutput func(stdout, stderr []byte), cancel chan bool) (*otto.MessageActionResult, *Error) {
	conn, err := host.connect()
	if err != nil {
		heartbeatStore.UpdateHostReachability(host, false)
		log.Error("Error triggering action on host '%s': %s", host.ID, err.Error())
		return nil, ErrorFrom(err)
	}
	defer conn.Close()

	result, err := conn.Conn.TriggerAction(action, actionOutput, cancel)
	if err != nil {
		return nil, ErrorUser(err.Error())
	}
	heartbeatStore.UpdateHostReachability(host, true)

	return result, nil
}

// Ping ping the host
func (host *Host) Ping() *Error {
	conn, err := host.connect()
	if err != nil {
		log.Error("Error sending heartbeat request to host '%s': %s", host.ID, err.Error())
		heartbeatStore.UpdateHostReachability(host, false)
		return ErrorFrom(err)
	}
	defer conn.Close()

	reply, err := conn.Conn.SendHeartbeat(otto.MessageHeartbeatRequest{ServerVersion: ServerVersion})
	if err != nil {
		log.PError("Error sending heartbeat request to host", map[string]interface{}{
			"host_id": host.ID,
			"error":   err.Error(),
		})
		heartbeatStore.UpdateHostReachability(host, false)
		return ErrorFrom(err)
	}
	heartbeatStore.RegisterHeartbeatReply(host, *reply)

	return nil
}

// RunScript run the script on the host. Error will only ever be populated with internal server
// errors, such as being unable to read from the database.
func (host *Host) RunScript(script *Script, scriptOutput func(stdout, stderr []byte), cancel chan bool) (*ScriptResult, *Error) {
	start := time.Now()

	sr, err := script.OttoScript()
	if err != nil {
		return nil, err
	}
	scriptRequest := *sr

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
		return nil, err
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

	scriptRequest.Environment = environ.Map(variables)

	log.Info("Executing script '%s' on host '%s'", script.Name, host.Address)
	result, err := host.TriggerAction(otto.MessageTriggerAction{
		Action: otto.ActionRunScript,
		Script: scriptRequest,
	}, scriptOutput, cancel)
	if result == nil && err == nil {
		err = ErrorServer("Unexpected end of connection")
	}
	if err != nil {
		log.Error("Error running script on host '%s': %s", host.Address, err.Message)
		return &ScriptResult{
			ScriptID:    script.ID,
			Duration:    time.Since(start),
			Environment: variables,
			Result: otto.ScriptResult{
				Success: false,
			},
			RunError: err.Message,
		}, nil
	}

	if result.ScriptResult.Success {
		if script.AfterExecution != "" {
			clientAction, ok := clientActionMap[script.AfterExecution]
			if !ok {
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

			log.Info("Performing post-execution action '%s' on host '%s'", script.AfterExecution, host.Address)
			_, err = host.TriggerAction(otto.MessageTriggerAction{
				Action: clientAction,
			}, nil, cancel)
			if err != nil {
				log.Error("Error running post-execution from script '%s' on host '%s': %s", script.Name, host.Address, result.ScriptResult.ExecError)
				return &ScriptResult{
					ScriptID:    script.ID,
					Duration:    time.Since(start),
					Environment: variables,
					Result: otto.ScriptResult{
						Success: false,
					},
					RunError: err.Message,
				}, nil
			}
		}
	} else {
		log.Error("Error running script '%s' on host '%s': %s", script.Name, host.Address, result.ScriptResult.ExecError)
	}

	return &ScriptResult{
		ScriptID:    script.ID,
		Duration:    time.Since(start),
		Environment: variables,
		Result:      result.ScriptResult,
	}, nil
}

// ExitClient exit the otto client on the host
func (host *Host) ExitClient() *Error {
	_, err := host.TriggerAction(otto.MessageTriggerAction{
		Action: otto.ActionExitClient,
	}, nil, nil)
	if err != nil {
		log.Error("Error exiting otto client on host '%s': %s", host.Address, err.Message)
		return err
	}
	return nil
}

func (host *Host) RotateIdentity() (string, string, *Error) {
	serverId, iderr := otto.NewIdentity()
	if iderr != nil {
		log.PError("Error generating new identity", map[string]interface{}{
			"error": iderr.Error(),
		})
		return "", "", ErrorFrom(iderr)
	}

	result, err := host.TriggerAction(otto.MessageTriggerAction{
		Action: otto.ActionUpdateIdentity,
		File: otto.File{
			Data: []byte(serverId.PublicKeyString()),
		},
	}, nil, nil)
	if err != nil {
		log.PError("Error requesting client update identity", map[string]interface{}{
			"host":  host.ID,
			"error": err.Message,
		})
		return "", "", err
	}
	if result.File.Data == nil {
		log.PError("Client did not return new identity", map[string]interface{}{
			"host":  host.ID,
			"error": err.Message,
		})
		return "", "", err
	}

	IdentityStore.Set(host.ID, serverId)
	err = HostStore.UpdateHostTrust(host.ID, HostTrust{
		TrustedIdentity: string(result.File.Data),
		LastTrustUpdate: time.Now(),
	})
	if err != nil {
		log.PError("Error updating host trust", map[string]interface{}{
			"host":  host.ID,
			"error": err.Message,
		})
		return "", "", err
	}

	log.PInfo("Rotated host identities", map[string]interface{}{
		"host_id":    host.ID,
		"host_name":  host.Name,
		"server_pub": serverId.PublicKeyString(),
		"client_pub": string(result.File.Data),
	})
	return serverId.PublicKeyString(), string(result.File.Data), nil
}