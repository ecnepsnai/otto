package server

import (
	"fmt"
	"io"
	"net"
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
	ClientActionExitClient: otto.ActionExit,
}

type hostConnection struct {
	Host    *Host
	Address string
	c       net.Conn
}

// connect will open a connection to the Otto client on the host
func (host *Host) connect() (*hostConnection, error) {
	address := fmt.Sprintf("%s:%d", host.Address, host.Port)
	log.Debug("Connecting to %s...", address)

	timeout := time.Duration(Options.Network.Timeout) * time.Second
	network := "tcp"
	if Options.Network.ForceIPVersion == IPVersionOptionIPv4 {
		network = "tcp4"
	} else if Options.Network.ForceIPVersion == IPVersionOptionIPv6 {
		network = "tcp6"
	}

	c, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		heartbeatStore.MarkHostUnreachable(host)
		log.Error("Error connecting to host '%s': %s", address, err.Error())
		return nil, err
	}

	log.Debug("Connected!")
	return &hostConnection{
		Host:    host,
		Address: fmt.Sprintf("%s:%d", host.Address, host.Port),
		c:       c,
	}, nil
}

// SendMessage will send the given Otto message to the host
func (hc *hostConnection) SendMessage(messageType uint32, message interface{}) error {
	if err := otto.WriteMessage(messageType, message, hc.c, hc.Host.PSK); err != nil {
		log.Error("Error sending message to host '%s': %s", hc.Address, err.Error())
		return err
	}

	return nil
}

// ReadMessage will try and read a message from the Otto host.
func (hc *hostConnection) ReadMessage() (uint32, interface{}, error) {
	return otto.ReadMessage(hc.c, hc.Host.PSK)
}

// Close will close the connection to the Otto host
func (hc *hostConnection) Close() {
	hc.c.Close()
}

// TriggerAction will trigger the given action on the host
func (host *Host) TriggerAction(action otto.MessageTriggerAction, actionOutput func(stdout, stderr []byte), cancel chan bool) (*otto.ScriptResult, *Error) {
	conn, err := host.connect()
	if err != nil {
		log.Error("Error triggering action on host '%s': %s", host.ID, err.Error())
		return nil, ErrorFrom(err)
	}
	defer conn.Close()

	if err := conn.SendMessage(otto.MessageTypeTriggerAction, action); err != nil {
		log.Error("Error triggering action on host '%s': %s", host.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	go func() {
		for {
			select {
			case <-cancel:
				conn.SendMessage(otto.MessageTypeCancelAction, otto.MessageCancelAction{})
			}
		}
	}()

	for true {
		messageType, message, err := conn.ReadMessage()
		if err == io.EOF || messageType == 0 {
			return nil, nil
		} else if err != nil {
			log.Error("Error triggering action on host '%s': %s", host.ID, err.Error())
			return nil, ErrorFrom(err)
		}

		switch messageType {
		case otto.MessageTypeActionOutput:
			output := message.(otto.MessageActionOutput)
			if actionOutput != nil {
				actionOutput(output.Stdout, output.Stderr)
			}
			break
		case otto.MessageTypeActionResult:
			result := message.(otto.MessageActionResult)
			scriptResult := &result.ScriptResult
			heartbeatStore.MarkHostReachable(host, result.ClientVersion)
			log.Debug("Action completed with result: %+v", *scriptResult)
			return scriptResult, nil
		case otto.MessageTypeGeneralFailure:
			result := message.(otto.MessageGeneralFailure)
			generalError := result.Error
			log.Error("General error triggering action on host '%s': %s", host.ID, generalError.Error())
			return nil, ErrorUser(generalError.Error())
		}
	}

	return nil, nil
}

// Ping ping the host
func (host *Host) Ping() *Error {
	conn, err := host.connect()
	if err != nil {
		log.Error("Error sending heartbeat request to host '%s': %s", host.ID, err.Error())
		return ErrorFrom(err)
	}
	defer conn.Close()

	if err := conn.SendMessage(otto.MessageTypeHeartbeatRequest, otto.MessageHeartbeatRequest{ServerVersion: ServerVersion}); err != nil {
		log.Error("Error sending heartbeat request to host '%s': %s", host.ID, err.Error())
		return ErrorFrom(err)
	}
	messageType, message, err := conn.ReadMessage()
	if err == io.EOF {
		log.Error("Client closed connection before replying to heartbeat '%s'", host.ID)
		return ErrorServer("Client closed connection")
	} else if err != nil {
		log.Error("Error sending heartbeat request to host '%s': %s", host.ID, err.Error())
		return ErrorFrom(err)
	}
	switch messageType {
	case otto.MessageTypeHeartbeatResponse:
		response := message.(otto.MessageHeartbeatResponse)
		heartbeatStore.MarkHostReachable(host, response.ClientVersion)
		break
	default:
		log.Error("Unexpected otto message %d while looking for heartbeat reply", messageType)
		return ErrorServer("Unexpected response")
	}

	return nil
}

// RunScript run the script on the host. Error will only ever be populated with internal server
// errors, such as being unable to read from the database.
func (host *Host) RunScript(script *Script, scriptOutput func(stdout, stderr []byte), cancel chan bool) (*ScriptResult, *Error) {
	start := time.Now()

	fileIDs, err := script.Attachments()
	if err != nil {
		return nil, err
	}
	files := make([]otto.File, len(script.AttachmentIDs))
	for i, file := range fileIDs {
		file, erro := file.OttoFile()
		if erro != nil {
			return nil, ErrorFrom(erro)
		}
		files[i] = *file
	}

	scriptRequest := otto.Script{
		Name:             script.Name,
		UID:              script.UID,
		GID:              script.GID,
		Executable:       script.Executable,
		Data:             []byte(script.Script),
		WorkingDirectory: script.WorkingDirectory,
		Environment:      map[string]string{},
		Files:            files,
	}

	variables := environ.Merge(staticEnvironment(), []environ.Variable{
		environ.New("OTTO_HOST_ADDRESS", host.Address),
		environ.New("OTTO_HOST_PORT", fmt.Sprintf("%d", host.Port)),
	})

	if Options.Security.IncludePSKEnv {
		variables = append(variables, environ.Variable{Key: "OTTO_HOST_PSK", Value: host.PSK, Secret: true})
	}

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

	if result.Success {
		if script.AfterExecution != "" {
			log.Info("Performing post-execution action '%s' on host '%s'", script.AfterExecution, host.Address)
			_, err = host.TriggerAction(otto.MessageTriggerAction{
				Action: clientActionMap[script.AfterExecution],
			}, nil, cancel)
			if err != nil {
				log.Error("Error running post-execution from script '%s' on host '%s': %s", script.Name, host.Address, result.ExecError)
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
		log.Error("Error running script '%s' on host '%s': %s", script.Name, host.Address, result.ExecError)
	}

	return &ScriptResult{
		ScriptID:    script.ID,
		Duration:    time.Since(start),
		Environment: variables,
		Result:      *result,
	}, nil
}

// ExitClient exit the otto client on the host
func (host *Host) ExitClient() *Error {
	_, err := host.TriggerAction(otto.MessageTriggerAction{
		Action: otto.ActionExit,
	}, nil, nil)
	if err != nil {
		log.Error("Error exiting otto client on host '%s': %s", host.Address, err.Message)
		return err
	}
	return nil
}
