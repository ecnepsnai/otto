package server

import (
	"fmt"
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

// PerformRequest perform an otto request on a host
func (host *Host) PerformRequest(request otto.Request) (*otto.Reply, *Error) {
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
		return nil, ErrorFrom(err)
	}
	log.Debug("Connected!")
	defer c.Close()

	if err := otto.WriteRequest(request, host.PSK, c); err != nil {
		log.Error("Error writing request to host '%s': %s", address, err.Error())
		return nil, ErrorFrom(err)
	}

	reply, err := otto.ReadReply(c, host.PSK)
	if err != nil {
		log.Error("Error reading reply from host '%s': %s", address, err.Error())
		return nil, ErrorFrom(err)
	}

	heartbeatStore.MarkHostReachable(host, reply)
	return reply, nil
}

// Ping ping the host
func (host *Host) Ping() *Error {
	_, err := host.PerformRequest(otto.Request{
		Action: otto.ActionPing,
	})
	if err != nil {
		log.Error("Error pinging host '%s': %s", host.Address, err.Message)
		return err
	}
	return nil
}

// RunScript run the script on the host. Error will only ever be populated with internal server
// errors, such as being unable to read from the database.
func (host *Host) RunScript(script *Script) (*ScriptResult, *Error) {
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
	reply, err := host.PerformRequest(otto.Request{
		Action: otto.ActionRunScript,
		Script: scriptRequest,
	})
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

	result := reply.ScriptResult
	if result.Success {
		log.Info("Result: OK")
		if script.AfterExecution != "" {
			log.Info("Performing post-execution action '%s' on host '%s'", script.AfterExecution, host.Address)
			_, err = host.PerformRequest(otto.Request{
				Action: clientActionMap[script.AfterExecution],
			})
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
		Result:      result,
	}, nil
}

// ExitClient exit the otto client on the host
func (host *Host) ExitClient() *Error {
	_, err := host.PerformRequest(otto.Request{
		Action: otto.ActionExit,
	})
	if err != nil {
		log.Error("Error exiting otto client on host '%s': %s", host.Address, err.Message)
		return err
	}
	return nil
}
