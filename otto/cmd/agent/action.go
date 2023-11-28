package main

import (
	"os"
	"os/exec"

	"github.com/ecnepsnai/otto/shared/otto"
)

func handleTriggerAction(conn *otto.Connection, messageType otto.MessageType, message interface{}) {
	switch messageType {
	case otto.MessageTypeTriggerActionRunScript:
		handleTriggerActionRunScript(conn, message.(otto.MessageTriggerActionRunScript))
	case otto.MessageTypeTriggerActionReloadConfig:
		conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{
			Error: handleTriggerActionReloadConfig(conn),
		})
	case otto.MessageTypeTriggerActionUploadFile:
		conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{
			Error: handleTriggerActionUploadFile(conn, message.(otto.MessageTriggerActionUploadFile)),
		})
	case otto.MessageTypeTriggerActionExitAgent:
		go handleTriggerActionExitAgent(conn)
		conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{})
	case otto.MessageTypeTriggerActionReboot:
		go handleTriggerActionReboot(conn)
		conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{})
	case otto.MessageTypeTriggerActionShutdown:
		go handleTriggerActionShutdown(conn)
		conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{})
	}
}

func handleTriggerActionReloadConfig(conn *otto.Connection) string {
	if err := loadConfig(); err != nil {
		return err.Error()
	}

	restartServer = true
	mustLoadIdentity()
	conn.Close()
	defer listener.Close()
	return ""
}

func handleTriggerActionExitAgent(conn *otto.Connection) {
	conn.Close()
	log.Warn("Exiting at request of '%s'", conn.RemoteAddr().String())
	os.Exit(1)
}

func handleTriggerActionReboot(conn *otto.Connection) {
	conn.Close()
	log.Warn("Rebooting at request of '%s'", conn.RemoteAddr().String())
	command := "/usr/sbin/reboot"
	if config.RebootCommand != nil {
		command = *config.RebootCommand
	}
	exec.Command(command).Run()
}

func handleTriggerActionShutdown(conn *otto.Connection) {
	conn.Close()
	log.Warn("Shutting down at request of '%s'", conn.RemoteAddr().String())
	command := "/usr/sbin/halt"
	if config.ShutdownCommand != nil {
		command = *config.ShutdownCommand
	}
	exec.Command(command).Run()
}
