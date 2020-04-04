package server

import (
	"fmt"
	"net"
	"time"

	"github.com/ecnepsnai/otto"
)

// ScriptResult describes a script result
type ScriptResult struct {
	ScriptID    string
	Duration    time.Duration
	Environment map[string]string
	Result      otto.ScriptResult
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

	heartbeatStore.MarkHostReachable(host)
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

// RunScript run the script on the host
func (host *Host) RunScript(script *Script) (*ScriptResult, *Error) {
	start := time.Now()

	scriptRequest := otto.Script{
		Name:             script.Name,
		UID:              script.UID,
		GID:              script.GID,
		Executable:       script.Executable,
		Data:             []byte(script.Script),
		WorkingDirectory: script.WorkingDirectory,
		Environment:      map[string]string{},
	}

	for key, val := range staticEnvironment() {
		scriptRequest.Environment[key] = val
	}
	scriptRequest.Environment["OTTO_HOST_ADDRESS"] = host.Address
	scriptRequest.Environment["OTTO_HOST_PORT"] = fmt.Sprintf("%d", host.Port)
	scriptRequest.Environment["OTTO_HOST_PSK"] = host.PSK

	// 1. Global envrionment variables
	for key, val := range Options.GlobalEnvironment {
		scriptRequest.Environment[key] = val
	}

	// 2. Script envrionment variables
	for key, val := range script.Environment {
		scriptRequest.Environment[key] = val
	}

	// 3. Group envrionment variables
	groups, err := host.Groups()
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		for key, val := range group.Environment {
			scriptRequest.Environment[key] = val
		}
	}

	// 4. Host envrionment variables
	for key, val := range host.Environment {
		scriptRequest.Environment[key] = val
	}

	log.Info("Executing script '%s' on host '%s'", script.Name, host.Address)
	reply, err := host.PerformRequest(otto.Request{
		Action: otto.ActionRunScript,
		Script: scriptRequest,
	})
	if err != nil {
		log.Error("Error running script on host '%s': %s", host.Address, err.Message)
		return nil, err
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
				return nil, err
			}
		}
	} else {
		log.Error("Error running script '%s' on host '%s': %s", script.Name, host.Address, result.ExecError)
	}

	return &ScriptResult{
		ScriptID:    script.ID,
		Duration:    time.Since(start),
		Environment: scriptRequest.Environment,
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
