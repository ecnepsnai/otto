package main

import (
	"encoding/base64"
	"io"
	"os/user"
	"strconv"
	"sync"
	"time"

	"github.com/ecnepsnai/otto/shared/otto"
)

var restartServer = false
var listener *otto.Listener
var identityLock = &sync.RWMutex{}

func listen() {
	for {
		var err error
		listener, err = otto.SetupListener(&otto.ListenOptions{
			Address:   config.ListenAddr,
			AllowFrom: getAllowFroms(),
			Identity:  agentIdentity,
			GetTrustedPublicKeys: func() []string {
				identityLock.RLock()
				defer identityLock.RUnlock()
				return []string{config.ServerIdentity}
			},
		}, handle)
		if err != nil {
			panic("error listening: " + err.Error())
		}
		listener.Accept()
		log.Warn("Listener stopped")
		if !restartServer {
			break
		}
		log.Info("Listener restarting")
	}
}

func handle(conn *otto.Connection) {
	log.PInfo("Connection established", map[string]interface{}{
		"remote_addr": conn.RemoteAddr().String(),
		"identity":    base64.StdEncoding.EncodeToString(conn.RemoteIdentity()),
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
		case otto.MessageTypeKeepalive:
			// Noop
		case otto.MessageTypeHeartbeatRequest:
			handleHeartbeatRequest(conn, message.(otto.MessageHeartbeatRequest))
		case otto.MessageTypeRotateIdentityRequest:
			handleRotateIdentity(conn, message.(otto.MessageRotateIdentityRequest))
		case otto.MessageTypeTriggerActionRunScript,
			otto.MessageTypeTriggerActionReloadConfig,
			otto.MessageTypeTriggerActionUploadFile,
			otto.MessageTypeTriggerActionExitAgent,
			otto.MessageTypeTriggerActionReboot,
			otto.MessageTypeTriggerActionShutdown:
			handleTriggerAction(conn, messageType, message, cancel)
		case otto.MessageTypeCancelAction:
			cancel <- true
		default:
			log.Warn("Unexpected message with type %d from %s", messageType, conn.RemoteAddr().String())
		}
	}
}

func handleHeartbeatRequest(conn *otto.Connection, message otto.MessageHeartbeatRequest) {
	log.PInfo("Incoming heartbeat", map[string]interface{}{
		"remote_addr":    conn.RemoteAddr().String(),
		"server_version": message.Version,
		"nonce":          message.Nonce,
	})
	Stats.LastHeartbeat = time.Now().UTC().Unix()

	properties := map[string]string{
		"hostname":             registerProperties.Hostname,
		"kernel_name":          registerProperties.KernelName,
		"kernel_version":       registerProperties.KernelVersion,
		"distribution_name":    registerProperties.DistributionName,
		"distribution_version": registerProperties.DistributionVersion,
	}

	response := otto.MessageHeartbeatResponse{
		AgentVersion: Version,
		Properties:   properties,
		Nonce:        message.Nonce,
	}

	if err := conn.WriteMessage(otto.MessageTypeHeartbeatResponse, response); err != nil {
		log.Error("Error writing heartbeat response to '%s': %s", conn.RemoteAddr().String(), err.Error())
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
	newID, err := loadAgentIdentity()
	if err != nil {
		reply.Error = err.Error()
	}

	newPublicKey := base64.StdEncoding.EncodeToString(newID.PublicKey().Marshal())
	reply.PublicKey = newPublicKey
	log.PWarn("Identity rotated", map[string]interface{}{
		"agent_public":  newPublicKey,
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
