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

func listen() {
	_, network, err := net.ParseCIDR(config.AllowFrom)
	if err != nil {
		panic("invalid CIDR address in property allow_from")
	}

	otto.Listen(otto.ListenOptions{
		Address:          config.ListenAddr,
		AllowFrom:        network,
		Identity:         clientIdentity,
		TrustedPublicKey: config.ServerIdentity,
	}, func(c *otto.Connection) {
		handle(c)
	})
}

func handle(conn *otto.Connection) {
	log.PInfo("Connection established", map[string]interface{}{
		"remote_addr": conn.RemoteAddr().String(),
	})
	defer conn.Close()

	cancel := make(chan bool)

	for {
		messageType, message, err := conn.ReadMessage()
		if err == io.EOF || messageType == 0 {
			break
		}
		if err != nil {
			log.Error("Error reading message from server '%s': %s", conn.RemoteAddr().String(), err.Error())
			break
		}
		log.Debug("Message from %s: %d", conn.RemoteAddr().String(), messageType)

		switch messageType {
		case otto.MessageTypeHeartbeatRequest:
			handleHeartbeatRequest(conn, message.(otto.MessageHeartbeatRequest))
		case otto.MessageTypeTriggerAction:
			go handleTriggerAction(conn, message.(otto.MessageTriggerAction), cancel)
		case otto.MessageTypeCancelAction:
			cancel <- true
		default:
			log.Warn("Unexpected message with type %d from %s", messageType, conn.RemoteAddr().String())
		}
	}
}

func handleHeartbeatRequest(conn *otto.Connection, message otto.MessageHeartbeatRequest) {
	log.Info("Heartbeat from %s (v%s)", conn.RemoteAddr().String(), message.ServerVersion)

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

	if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, response); err != nil {
		log.Error("Error writing heartbeat response to '%s': %s", conn.RemoteAddr().String(), err.Error())
	}
}

func handleTriggerAction(conn *otto.Connection, message otto.MessageTriggerAction, cancel chan bool) {
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
		reply.ScriptResult = runScript(conn, message.Script, cancel)
	case otto.ActionUploadFile, otto.ActionUploadFileAndExit:
		if err := uploadFile(message.File); err != nil {
			reply.Error = err
		}
	case otto.ActionUpdateIdentity:
		newIdentity := string(message.File.Data)
		if newIdentity == "" {
			reply.Error = fmt.Errorf("no identity provided")
		} else {
			if err := updateServerIdentity(newIdentity); err != nil {
				reply.Error = err
			}
		}
	default:
		log.Error("Unknown action %d", message.Action)
		return
	}

	log.Debug("Trigger complete, writing reply")
	if err := conn.WriteMessage(otto.MessageTypeActionResult, reply); err != nil {
		log.Error("Error writing reply to '%s': %s", conn.RemoteAddr().String(), err.Error())
		return
	}

	if message.Action == otto.ActionUpdateIdentity {
		loadConfig()
	}

	if message.Action == otto.ActionUploadFileAndExit || message.Action == otto.ActionExit {
		conn.Close()
		log.Warn("Exiting at request of '%s'", conn.RemoteAddr().String())
		os.Exit(1)
	}

	if message.Action == otto.ActionReboot {
		conn.Close()
		log.Warn("Rebooting at request of '%s'", conn.RemoteAddr().String())
		exec.Command("/usr/sbin/reboot").Run()
	} else if message.Action == otto.ActionShutdown {
		conn.Close()
		log.Warn("Shutting down at request of '%s'", conn.RemoteAddr().String())
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
