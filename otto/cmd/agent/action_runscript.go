package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/ecnepsnai/otto/shared/otto"
)

var scriptLog = &sync.Map{}

type scriptOutputWriter struct {
	isStderr  bool
	conn      *otto.Connection
	file      *os.File
	lastWrite time.Time
}

func (w *scriptOutputWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if err := w.conn.WriteMessage(otto.MessageTypeActionOutput, otto.MessageActionOutput{
		IsStdErr: w.isStderr,
		Data:     p,
	}); err != nil {
		return n, err
	}
	if _, err := w.file.Write(p); err != nil {
		return n, err
	}
	w.lastWrite = time.Now()
	return n, nil
}

func handleTriggerActionRunScript(conn *otto.Connection, message otto.MessageTriggerActionRunScript) {
	if pid, running := scriptLog.Load(message.Name); running {
		log.Error("Cannot start script '%s' as it's already running on pid %d", message.Name, pid.(int))
		conn.WriteMessage(otto.MessageTypeActionResult, otto.ScriptResult{
			Success:   false,
			ExecError: fmt.Sprintf("script '%s' already running on host", message.Name),
		})
		conn.Close()
		return
	}

	start := time.Now()

	tmpDir, err := os.MkdirTemp("", "otto")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	scriptPath := path.Join(tmpDir, "script")
	scriptFileWriter, err := os.OpenFile(scriptPath, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		panic(err)
	}
	log.Debug("Writing script to %s", scriptPath)

	stdoutPath := path.Join(tmpDir, "stdout")
	stderrPath := path.Join(tmpDir, "stderr")
	combinedPath := path.Join(tmpDir, "combined")
	stdoutFile, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	stderrFile, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	combinedFile, err := os.OpenFile(combinedPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	stdout := &scriptOutputWriter{
		conn:      conn,
		file:      stdoutFile,
		lastWrite: time.Now(),
	}
	stderr := &scriptOutputWriter{
		isStderr:  true,
		conn:      conn,
		file:      stderrFile,
		lastWrite: time.Now(),
	}

	totalScriptWrote := uint64(0)
	var scriptBuffer = make([]byte, min(message.ScriptInfo.Length, 1024))
	log.Debug("Allocating %dB for script buffer", len(scriptBuffer))
	readAttempts := 0

	log.Debug("Telling server we're ready for script data")
	if err := conn.WriteMessage(otto.MessageTypeReadyForData, nil); err != nil {
		log.PError("Error replying to server", map[string]interface{}{
			"error": err.Error(),
		})
		conn.WriteMessage(otto.MessageTypeActionResult, otto.ScriptResult{
			Success:   false,
			ExecError: "Error copying script data",
			Elapsed:   time.Since(start),
		})
		conn.Close()
		return
	}

	for totalScriptWrote < message.ScriptInfo.Length {
		read, err := conn.ReadData(scriptBuffer)
		if err != nil && err != io.EOF {
			scriptFileWriter.Close()
			log.PError("Error reading script data", map[string]interface{}{
				"error": err.Error(),
			})
			conn.WriteMessage(otto.MessageTypeActionResult, otto.ScriptResult{
				Success:   false,
				ExecError: "Error copying script data",
				Elapsed:   time.Since(start),
			})
			conn.Close()
			return
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

		wrote, err := scriptFileWriter.Write(scriptBuffer[0:read])
		if err != nil {
			scriptFileWriter.Close()
			log.PError("Error writing script data", map[string]interface{}{
				"error": err.Error(),
			})
			conn.WriteMessage(otto.MessageTypeActionResult, otto.ScriptResult{
				Success:   false,
				ExecError: "Error copying script data",
				Elapsed:   time.Since(start),
			})
			conn.Close()
			return
		}
		totalScriptWrote += uint64(read)
		log.Debug("Wrote %dB of script data to %s", wrote, scriptPath)
	}
	scriptFileWriter.Close()

	if totalScriptWrote != message.ScriptInfo.Length {
		log.PError("Unexpected end of script data", map[string]interface{}{
			"script_length": message.ScriptInfo.Length,
			"wrote_length":  totalScriptWrote,
		})
		conn.WriteMessage(otto.MessageTypeActionResult, otto.ScriptResult{
			Success:   false,
			ExecError: "Error copying script data",
			Elapsed:   time.Since(start),
		})
		conn.Close()
		return
	}

	log.PInfo("Executing script", map[string]interface{}{
		"remote_addr": conn.RemoteAddr().String(),
		"name":        message.Name,
		"wd":          message.WorkingDirectory,
		"exec":        message.Executable,
		"environ":     message.Environment,
	})
	Stats.ScriptsExecuted++
	Stats.LastScriptExecuted = time.Now().UTC().Unix()

	cmd := exec.Command(message.Executable, scriptPath)
	log.Debug("Exec: %s %s", message.Executable, scriptPath)

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
			conn.WriteMessage(otto.MessageTypeActionResult, otto.ScriptResult{
				Success:   false,
				ExecError: "Running a script as a specific user requires the Otto agent running as root",
				Elapsed:   time.Since(start),
			})
			conn.Close()
			return
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
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	log.Debug("Running '%s' %s", scriptPath, cmd.Args)
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
	isRunning := true
	proc := cmd.Process
	log.Debug("Waiting for script...")
	scriptLog.Store(message.Name, proc.Pid)

	go func() {
		for isRunning {
			if time.Since(stdout.lastWrite) > 30*time.Second && time.Since(stderr.lastWrite) > 30*time.Second {
				conn.WriteMessage(otto.MessageTypeKeepalive, nil)
				stdout.lastWrite = time.Now()
				stderr.lastWrite = time.Now()
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	cmd.Wait()
	isRunning = false
	scriptLog.Delete(message.Name)

	stdoutFile.Sync()
	stderrFile.Sync()
	stdoutFile.Seek(0, 0)
	stderrFile.Seek(0, 0)
	stdoutLen, _ := io.Copy(combinedFile, stdoutFile)
	combinedFile.Sync()
	stdoutFile.Close()
	os.Remove(stdoutPath)
	stderrLen, _ := io.Copy(combinedFile, stderrFile)
	stderrFile.Close()
	combinedFile.Sync()
	combinedFile.Seek(0, 0)
	os.Remove(stderrPath)

	log.PInfo("Finished executing script", map[string]interface{}{
		"remote_addr": conn.RemoteAddr().String(),
		"elapsed":     time.Since(start).String(),
		"exit_code":   cmd.ProcessState.ExitCode(),
		"name":        message.Name,
		"wd":          message.WorkingDirectory,
		"exec":        message.Executable,
	})

	result.StdoutLen = uint32(stdoutLen)
	result.StderrLen = uint32(stderrLen)
	result.Elapsed = time.Since(start)

	if cmd.ProcessState.ExitCode() != 0 {
		result.Success = false
		log.Error("Script exit code: %d", cmd.ProcessState.ExitCode())
	} else {
		result.Success = true
	}

	if err := conn.WriteMessage(otto.MessageTypeActionResult, otto.MessageActionResult{
		ScriptResult: result,
	}); err != nil {
		log.PError("Error writing action resuslt", map[string]interface{}{
			"error": err.Error(),
		})
		conn.Close()
		return
	}

	messageType, _, err := conn.ReadMessage()
	if err != nil {
		log.PError("Error waiting for reply from server", map[string]interface{}{
			"error": err.Error(),
		})
		conn.Close()
		return
	}
	if messageType != otto.MessageTypeReadyForData {
		log.Error("Unexpected message %d", messageType)
		conn.Close()
		return
	}
	if _, err := conn.Copy(combinedFile); err != nil {
		log.PError("Error writing script output", map[string]interface{}{
			"error": err.Error(),
		})
		conn.Close()
		return
	}
	conn.WriteFinished()
	log.Debug("Finished runscript action, closing connection...")
	conn.Close()
}

func handleCancelAction(conn *otto.Connection, message otto.MessageCancelAction) {
	pid, found := scriptLog.Load(message.Name)
	if !found {
		log.Warn("Attempt to cancel unknown script named '%s'", message.Name)
		return
	}

	log.Warn("Cancelling script '%s' on pid %d", message.Name, pid.(int))
	syscall.Kill(pid.(int), syscall.SIGKILL)
}
