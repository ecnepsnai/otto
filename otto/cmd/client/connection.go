package main

import (
	"encoding/base64"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"

	"github.com/ecnepsnai/otto"
)

var restartServer = false
var listener *otto.Listener

func listen() {
	for {
		var err error
		listener, err = otto.SetupListener(otto.ListenOptions{
			Address:          config.ListenAddr,
			AllowFrom:        getAllowFroms(),
			Identity:         clientIdentity,
			TrustedPublicKey: config.ServerIdentity,
		}, handle)
		if err != nil {
			panic("error listening: " + err.Error())
		}
		listener.Accept()
		log.Warn("Server stopped")
		if !restartServer {
			break
		}
		log.Info("Server restarting")
	}
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
		case otto.MessageTypeRotateIdentityRequest:
			handleRotateIdentity(conn, message.(otto.MessageRotateIdentityRequest))
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
	case otto.ActionReloadConfig, otto.ActionExitClient, otto.ActionReboot, otto.ActionShutdown:
		// No action
		break
	case otto.ActionRunScript:
		reply.ScriptResult = runScript(conn, message.Script, cancel)
	case otto.ActionUploadFile, otto.ActionUploadFileAndExitClient:
		if err := uploadFile(message.File); err != nil {
			reply.Error = err.Error()
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

	if message.Action == otto.ActionReloadConfig {
		if err := loadConfig(); err != nil {
			reply.Error = err.Error()
		}

		restartServer = true
		mustLoadIdentity()
		conn.Close()
		defer listener.Close()
		return
	}

	if message.Action == otto.ActionUploadFileAndExitClient || message.Action == otto.ActionExitClient {
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

func handleRotateIdentity(conn *otto.Connection, message otto.MessageRotateIdentityRequest) {
	reply := otto.MessageRotateIdentityResponse{}

	if message.PublicKey == "" {
		reply.Error = "no identity provided"
	} else {
		if err := updateServerIdentity(message.PublicKey); err != nil {
			reply.Error = err.Error()
		}
	}

	if err := generateIdentity(); err != nil {
		reply.Error = err.Error()
	}
	newID, err := loadClientIdentity()
	if err != nil {
		reply.Error = err.Error()
	}

	newPublicKey := base64.RawStdEncoding.EncodeToString(newID.PublicKey().Marshal())
	reply.PublicKey = newPublicKey
	log.PWarn("Identity rotated", map[string]interface{}{
		"client_public": newPublicKey,
		"server_public": message.PublicKey,
	})

	if err := conn.WriteMessage(otto.MessageTypeRotateIdentityResponse, reply); err != nil {
		log.PError("Error writing rotate identity response", map[string]interface{}{
			"remote_addr": conn.RemoteAddr().String(),
			"error":       err.Error(),
		})
	}

	if err := loadConfig(); err != nil {
		reply.Error = err.Error()
	}

	restartServer = true
	mustLoadIdentity()
	conn.Close()
	defer listener.Close()
	return
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
