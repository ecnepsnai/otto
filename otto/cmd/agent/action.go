package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/ecnepsnai/otto/shared/otto"
)

func handleTriggerAction(conn *otto.Connection, messageType otto.MessageType, message interface{}, cancel chan bool) {
	switch messageType {
	case otto.MessageTypeTriggerActionRunScript:
		finished := false
		startTime := time.Now().Unix()
		timeoutSeconds := int64(600)
		if config.ScriptTimeout != nil {
			timeoutSeconds = *config.ScriptTimeout
		}

		go func() {
			for !finished {
				time.Sleep(100 * time.Millisecond)
				if time.Now().Unix()-startTime > timeoutSeconds {
					cancel <- true
				}
			}
		}()

		result := handleTriggerActionRunScript(conn, message.(otto.MessageTriggerActionRunScript), cancel)
		conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{
			ScriptResult: result,
		})
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

func handleTriggerActionRunScript(conn *otto.Connection, message otto.MessageTriggerActionRunScript, cancel chan bool) otto.ScriptResult {
	var proc *os.Process
	canCancel := true
	go func() {
		for canCancel {
			select {
			case <-cancel:
				if proc != nil {
					pgid, err := syscall.Getpgid(proc.Pid)
					if err != nil {
						log.Error("Error trying to kill process: %s", err.Error())
					}
					syscall.Kill(-pgid, 15)
					log.Warn("Killed running script")
				}
				conn.Close()
			default:
				//
			}
		}
	}()

	start := time.Now()

	tmpDir, err := os.MkdirTemp("", "otto")
	if err != nil {
		panic(err)
	}

	tmp, err := os.CreateTemp("", "otto")
	if err != nil {
		panic(err)
	}
	log.Debug("Writing script to %s", tmp.Name())
	if err := tmp.Chmod(0700); err != nil {
		tmp.Close()
		panic(err)
	}
	defer os.Remove(tmp.Name())

	totalScriptWrote := uint64(0)
	var scriptBuffer = make([]byte, min(message.ScriptInfo.Length, 1024))
	log.Debug("Allocating %dB for script buffer", len(scriptBuffer))
	readAttempts := 0

	log.Debug("Telling server we're ready for script data")
	if err := conn.WriteMessage(otto.MessageTypeReadyForData, nil); err != nil {
		log.PError("Error replying to server", map[string]interface{}{
			"error": err.Error(),
		})
		canCancel = false
		return otto.ScriptResult{
			Success:   false,
			ExecError: "Error copying script data",
			Elapsed:   time.Since(start),
		}
	}

	for totalScriptWrote < message.ScriptInfo.Length {
		read, err := conn.ReadData(scriptBuffer)
		if err != nil && err != io.EOF {
			tmp.Close()
			log.PError("Error reading script data", map[string]interface{}{
				"error": err.Error(),
			})
			canCancel = false
			return otto.ScriptResult{
				Success:   false,
				ExecError: "Error copying script data",
				Elapsed:   time.Since(start),
			}
		}

		// If no data but we haven't finished copying yet, give the server a bit more time
		if read == 0 {
			if totalScriptWrote < message.ScriptInfo.Length && readAttempts < 15 {
				log.Debug("No script data yet, waiting for server (%d/15)", readAttempts+1)
				time.Sleep(50 * time.Millisecond)
				readAttempts++
				continue
			}

			log.Debug("Finished reading script")
			break
		}

		wrote, err := tmp.Write(scriptBuffer[0:read])
		if err != nil {
			tmp.Close()
			log.PError("Error writing script data", map[string]interface{}{
				"error": err.Error(),
			})
			canCancel = false
			return otto.ScriptResult{
				Success:   false,
				ExecError: "Error copying script data",
				Elapsed:   time.Since(start),
			}
		}
		totalScriptWrote += uint64(read)
		log.Debug("Wrote %dB of script data to %s", wrote, tmp.Name())
	}
	tmp.Close()

	if totalScriptWrote != message.ScriptInfo.Length {
		log.PError("Unexpected end of script data", map[string]interface{}{
			"script_length": message.ScriptInfo.Length,
			"wrote_length":  totalScriptWrote,
		})
		canCancel = false
		return otto.ScriptResult{
			Success:   false,
			ExecError: "Error copying script data",
			Elapsed:   time.Since(start),
		}
	}

	log.PInfo("Executing script", map[string]interface{}{
		"remote_addr": conn.RemoteAddr().String(),
		"name":        message.Name,
		"wd":          message.WorkingDirectory,
		"exec":        message.Executable,
	})
	Stats.ScriptsExecuted++
	Stats.LastScriptExecuted = time.Now().UTC().Unix()

	cmd := exec.Command(message.Executable, tmp.Name())
	log.Debug("Exec: %s %s", message.Executable, tmp.Name())

	if len(message.Environment) > 0 {
		env := make([]string, len(message.Environment))
		i := 0
		for key, val := range message.Environment {
			env[i] = key + "=" + val
			i++
		}
		cmd.Env = env
		log.Debug("Environment: %v", env)
	}
	if config.Path != "" {
		cmd.Env = append(cmd.Env, "PATH="+config.Path)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	if !message.RunAs.Inherit {
		if uid, _ := getCurrentUIDandGID(); uid != 0 {
			log.Error("Cannot run script as specific user without the Otto agent running as root")
			canCancel = false
			return otto.ScriptResult{
				Success:   false,
				ExecError: "Running a script as a specific user requires the Otto agent running as root",
				Elapsed:   time.Since(start),
			}
		}

		log.Debug("Using UID %d and GID %d\n", message.RunAs.UID, message.RunAs.GID)
		cmd.SysProcAttr.Credential = &syscall.Credential{
			Uid: message.RunAs.UID,
			Gid: message.RunAs.GID,
		}
	}

	if message.WorkingDirectory != "" {
		cmd.Dir = message.WorkingDirectory
	} else {
		cmd.Dir = tmpDir
	}

	result := otto.ScriptResult{}

	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	log.Debug("Running '%s'", tmp.Name())
	if err := cmd.Start(); err != nil {
		result.Success = false
		if exitError, ok := err.(*exec.ExitError); ok {
			result.Code = exitError.ExitCode()
			log.Error("Script exit code: %d", result.Code)
		} else {
			log.Error("Error running script: %s", err.Error())
			result.ExecError = err.Error()
		}
	}
	proc = cmd.Process
	log.Debug("Waiting for script...")
	didExit := false
	go func() {
		lastLen := 0
		lastKA := time.Now().AddDate(0, 0, -1)
		for !didExit {
			outputLength := stdout.Len() + stderr.Len()
			if outputLength > lastLen {
				lastLen = outputLength
				log.Debug("Read %dB from stdout & stderr", outputLength)
				conn.WriteMessage(otto.MessageTypeActionOutput, otto.MessageActionOutput{
					Stdout: stdout.Bytes(),
					Stderr: stderr.Bytes(),
				})
			}
			if time.Since(lastKA) > 10*time.Second {
				conn.WriteMessage(otto.MessageTypeKeepalive, nil)
				lastKA = time.Now()
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
	err = cmd.Wait()
	log.PInfo("Finished executing script", map[string]interface{}{
		"remote_addr": conn.RemoteAddr().String(),
		"elapsed":     time.Since(start).String(),
		"exit_code":   cmd.ProcessState.ExitCode(),
		"name":        message.Name,
		"wd":          message.WorkingDirectory,
		"exec":        message.Executable,
	})
	didExit = true

	result.Stderr = stderr.String()
	result.Stdout = stdout.String()
	result.Elapsed = time.Since(start)
	log.Debug("Stdout: %s", result.Stdout)
	log.Debug("Stderr: %s", result.Stderr)

	if err != nil {
		result.Success = false
		if exitError, ok := err.(*exec.ExitError); ok {
			result.Code = exitError.ExitCode()
			log.Error("Script exit code: %d", result.Code)
		} else {
			log.Error("Error running script: %s", err.Error())
			result.ExecError = err.Error()
		}
		canCancel = false
		return result
	}

	result.Success = true
	canCancel = false
	os.RemoveAll(tmpDir)
	return result
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

func handleTriggerActionUploadFile(conn *otto.Connection, message otto.MessageTriggerActionUploadFile) string {
	err := uploadFile(message.FileInfo, func(f io.Writer) error {
		log.Debug("Telling server we're ready for script data")
		if err := conn.WriteMessage(otto.MessageTypeReadyForData, nil); err != nil {
			log.PError("Error replying to server", map[string]interface{}{
				"error": err.Error(),
			})
			return err
		}

		totalCopied := uint64(0)
		var fileBuffer = make([]byte, 1024)
		for totalCopied < message.Length {
			read, err := conn.ReadData(fileBuffer)
			if err != nil && err != io.EOF {
				return err
			}
			f.Write(fileBuffer[0:read])
			totalCopied += uint64(read)
		}
		return nil
	})
	if err != nil {
		return err.Error()
	}
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
