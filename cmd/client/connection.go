package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strconv"

	"github.com/ecnepsnai/otto"
)

type serverConnection struct {
	c             net.Conn
	remoteAddress string
}

func newServerConnection(c net.Conn) *serverConnection {
	sc := serverConnection{
		c:             c,
		remoteAddress: c.RemoteAddr().String(),
	}
	return &sc
}

func (sc *serverConnection) Start() {
	log.Info("New connection from %s", sc.remoteAddress)
	defer sc.c.Close()

	cancel := make(chan bool)

	for {
		messageType, message, err := otto.ReadMessage(sc.c, config.PSK)
		if err == io.EOF || messageType == 0 {
			break
		}
		if err != nil {
			log.Error("Error reading message from server '%s': %s", sc.remoteAddress, err.Error())
			break
		}
		log.Debug("Message from %s: %d", sc.remoteAddress, messageType)

		switch messageType {
		case otto.MessageTypeHeartbeatRequest:
			handleHeartbeatRequest(sc.c, message.(otto.MessageHeartbeatRequest))
		case otto.MessageTypeTriggerAction:
			go handleTriggerAction(sc.c, message.(otto.MessageTriggerAction), cancel)
		case otto.MessageTypeCancelAction:
			cancel <- true
		default:
			log.Warn("Unexpected message with type %d from %s", messageType, sc.remoteAddress)
		}
	}
}

func handleHeartbeatRequest(c net.Conn, message otto.MessageHeartbeatRequest) {
	log.Info("Heartbeat from %s (v%s)", c.RemoteAddr().String(), message.ServerVersion)

	properties := map[string]string{
		"hostname":             registerProperties.Hostname,
		"kernel_name":          registerProperties.KernelName,
		"kernel_version":       registerProperties.KernelVersion,
		"distribution_name":    registerProperties.DistributionName,
		"distribution_version": registerProperties.DistributionVersion,
	}

	response := otto.MessageHeartbeatResponse{
		ClientVersion: MainVersion,
		Properties:    properties,
	}

	if err := otto.WriteMessage(otto.MessageTypeHeartbeatResponse, response, c, config.PSK); err != nil {
		log.Error("Error writing heartbeat response to '%s': %s", c.RemoteAddr().String(), err.Error())
	}
}

func handleTriggerAction(c net.Conn, message otto.MessageTriggerAction, cancel chan bool) {
	reply := otto.MessageActionResult{
		ClientVersion: MainVersion,
	}
	switch message.Action {
	case otto.ActionExit, otto.ActionReboot, otto.ActionShutdown:
		// No action
		break
	case otto.ActionReloadConfig:
		if err := loadConfig(); err != nil {
			reply.Error = err
		}
	case otto.ActionRunScript:
		reply.ScriptResult = runScript(c, message.Script, cancel)
	case otto.ActionUploadFile, otto.ActionUploadFileAndExit:
		if err := uploadFile(message.File); err != nil {
			reply.Error = err
		}
	case otto.ActionUpdatePSK:
		newPSK := message.NewPSK
		if newPSK == "" {
			reply.Error = fmt.Errorf("no psk provided")
		} else {
			if err := updatePSK(newPSK); err != nil {
				reply.Error = err
			}
		}
	default:
		log.Error("Unknown action %d", message.Action)
		return
	}

	log.Debug("Trigger complete, writing reply")
	if err := otto.WriteMessage(otto.MessageTypeActionResult, reply, c, config.PSK); err != nil {
		log.Error("Error writing reply to '%s': %s", c.RemoteAddr().String(), err.Error())
		return
	}

	if message.Action == otto.ActionUpdatePSK {
		loadConfig()
	}

	if message.Action == otto.ActionUploadFileAndExit || message.Action == otto.ActionExit {
		c.Close()
		log.Warn("Exiting at request of '%s'", c.RemoteAddr().String())
		os.Exit(1)
	}

	if message.Action == otto.ActionReboot {
		c.Close()
		log.Warn("Rebooting at request of '%s'", c.RemoteAddr().String())
		exec.Command("/usr/sbin/reboot").Run()
	} else if message.Action == otto.ActionShutdown {
		c.Close()
		log.Warn("Shutting down at request of '%s'", c.RemoteAddr().String())
		exec.Command("/usr/sbin/halt").Run()
	}
}

func getCurrentUIDandGID() (uint32, uint32) {
	current, err := user.Current()
	if err != nil {
		panic(err)
	}
	uid, err := strconv.ParseUint(current.Uid, 10, 32)
	if err != nil {
		panic(err)
	}
	gid, err := strconv.ParseUint(current.Gid, 10, 32)
	if err != nil {
		panic(err)
	}

	return uint32(uid), uint32(gid)
}
