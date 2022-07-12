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

func runScript(conn *otto.Connection, script otto.Script, cancel chan bool) otto.ScriptResult {
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
			default:
				//
			}
		}
	}()

	start := time.Now()
	log.PInfo("Executing script", map[string]interface{}{
		"remote_addr": conn.RemoteAddr().String(),
		"name":        script.Name,
		"wd":          script.WorkingDirectory,
		"exec":        script.Executable,
	})

	for _, file := range script.Files {
		if err := uploadFile(file); err != nil {
			log.Error("Error uploading script file '%s': %s", file.Path, err.Error())
			canCancel = false
			return otto.ScriptResult{
				Success:   false,
				ExecError: err.Error(),
				Elapsed:   time.Since(start),
			}
		}
	}

	tmp, err := os.CreateTemp("", "otto")
	if err != nil {
		panic(err)
	}
	log.Debug("Writing script to %s", tmp.Name())
	if err := tmp.Chmod(0777); err != nil {
		tmp.Close()
		panic(err)
	}
	defer os.Remove(tmp.Name())
	if _, err := io.CopyBuffer(tmp, bytes.NewBuffer(script.Data), nil); err != nil {
		tmp.Close()
		panic(err)
	}
	tmp.Close()
	cmd := exec.Command(script.Executable, tmp.Name())
	log.Debug("Exec: %s %s", script.Executable, tmp.Name())

	if len(script.Environment) > 0 {
		env := make([]string, len(script.Environment))
		i := 0
		for key, val := range script.Environment {
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
	if !script.RunAs.Inherit {
		if uid, _ := getCurrentUIDandGID(); uid != 0 {
			log.Error("Cannot run script as specific user without the Otto agent running as root")
			canCancel = false
			return otto.ScriptResult{
				Success:   false,
				ExecError: "Running a script as a specific user requires the Otto agent running as root",
				Elapsed:   time.Since(start),
			}
		}

		log.Debug("Using UID %d and GID %d\n", script.RunAs.UID, script.RunAs.GID)
		cmd.SysProcAttr.Credential = &syscall.Credential{
			Uid: script.RunAs.UID,
			Gid: script.RunAs.GID,
		}
	}

	if script.WorkingDirectory != "" {
		cmd.Dir = script.WorkingDirectory
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
		"name":        script.Name,
		"wd":          script.WorkingDirectory,
		"exec":        script.Executable,
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
	return result
}
