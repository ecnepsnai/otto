package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto"
)

var log *logtic.Source

func main() {
	parseArgs()
	tryAutoRegister()
	mustLoadConfig()

	logtic.Log.FilePath = path.Join(config.LogPath, "otto_client.log")
	logtic.Log.Level = logtic.LevelWarn
	env := envMap()
	if _, verbose := env["OTTO_VERBOSE"]; verbose {
		logtic.Log.Level = logtic.LevelDebug
	}

	logtic.Open()
	log = logtic.Connect("otto")

	l, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		panic(err)
	}
	log.Info("Otto client listening on %s", config.ListenAddr)
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}

		_, network, _ := net.ParseCIDR(config.AllowFrom)
		if !network.Contains(c.RemoteAddr().(*net.TCPAddr).IP) {
			log.Warn("Rejecting connection from server outside of allowed network: %s", c.RemoteAddr().String())
			c.Close()
			continue
		}

		c.RemoteAddr()
		go newMessage(c)
	}
}

func parseArgs() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}

	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-v" || arg == "--version" {
			fmt.Printf("Otto client %s, Protocol version: %d, Runtime %s\n", MainVersion, otto.ProtocolVersion, runtime.Version())
			os.Exit(0)
		}
		i++
	}
}

func newMessage(c net.Conn) {
	log.Info("New message from %s", c.RemoteAddr().String())
	defer c.Close()

	messageType, message, err := otto.ReadMessage(c, config.PSK)
	if err != nil {
		log.Error("Error reading message from server '%s': %s", c.RemoteAddr().String(), err.Error())
		return
	}
	log.Debug("Message from %s: %d", c.RemoteAddr().String(), messageType)

	switch messageType {
	case otto.MessageTypeHeartbeatRequest:
		handleHeartbeatRequest(c, message.(otto.MessageHeartbeatRequest))
		return
	case otto.MessageTypeTriggerAction:
		handleTriggerAction(c, message.(otto.MessageTriggerAction))
		return
	}

	log.Warn("Unexpected message with type %d from %s", messageType, c.RemoteAddr().String())
}

func handleHeartbeatRequest(c net.Conn, message otto.MessageHeartbeatRequest) {
	log.Info("Heartbeat from %s (v%s)", c.RemoteAddr().String(), message.ServerVersion)

	response := otto.MessageHeartbeatResponse{
		ClientVersion: MainVersion,
	}

	if err := otto.WriteMessage(otto.MessageTypeHeartbeatResponse, response, c, config.PSK); err != nil {
		log.Error("Error writing heartbeat response to '%s': %s", c.RemoteAddr().String(), err.Error())
	}
}

func handleTriggerAction(c net.Conn, message otto.MessageTriggerAction) {
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
		break
	case otto.ActionRunScript:
		reply.ScriptResult = runScript(c, message.Script)
		break
	case otto.ActionUploadFile, otto.ActionUploadFileAndExit:
		if err := uploadFile(message.File); err != nil {
			reply.Error = err
		}
		break
	default:
		log.Error("Unknown action %d", message.Action)
		return
	}

	log.Debug("Trigger complete, writing reply")
	if err := otto.WriteMessage(otto.MessageTypeActionResult, reply, c, config.PSK); err != nil {
		log.Error("Error writing reply to '%s': %s", c.RemoteAddr().String(), err.Error())
		return
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

func runScript(c net.Conn, script otto.Script) otto.ScriptResult {
	start := time.Now()

	for _, file := range script.Files {
		if err := uploadFile(file); err != nil {
			log.Error("Error uploading script file '%s': %s", file.Path, err.Error())
			return otto.ScriptResult{
				Success:   false,
				ExecError: err.Error(),
				Elapsed:   time.Since(start),
			}
		}
	}

	tmp, err := ioutil.TempFile("", "otto")
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

	if getCurrentUID() != script.UID || getCurrentGID() != script.GID {
		log.Debug("Using UID %d and GID %d\n", script.UID, script.GID)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: script.UID,
				Gid: script.GID,
			},
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

	log.Info("Running '%s' as UID %d GID %d", tmp.Name(), script.UID, script.GID)
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
	log.Debug("Waiting for script...")
	didExit := false
	go func() {
		lastLen := 0
		for !didExit {
			outputLength := stdout.Len() + stderr.Len()
			if outputLength > lastLen {
				lastLen = outputLength
				log.Debug("Read %dB from stdout & stderr", outputLength)
				otto.WriteMessage(otto.MessageTypeActionOutput, otto.MessageActionOutput{
					Stdout: stdout.Bytes(),
					Stderr: stderr.Bytes(),
				}, c, config.PSK)
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()
	err = cmd.Wait()
	log.Info("Script '%s' exited after %s", time.Since(start))
	didExit = true
	if err != nil {
		result.Success = false
		if exitError, ok := err.(*exec.ExitError); ok {
			result.Code = exitError.ExitCode()
			log.Error("Script exit code: %d", result.Code)
		} else {
			log.Error("Error running script: %s", err.Error())
			result.ExecError = err.Error()
		}
		return result
	}
	log.Info("Script exit OK")
	result.Success = true
	result.Stderr = string(stderr.Bytes())
	result.Stdout = string(stdout.Bytes())
	result.Elapsed = time.Since(start)
	log.Debug("Stdout: %s", result.Stdout)
	log.Debug("Stderr: %s", result.Stderr)

	return result
}

func uploadFile(file otto.File) error {
	f, err := os.OpenFile(file.Path, os.O_CREATE|os.O_WRONLY, os.FileMode(file.Mode))
	if err != nil {
		log.Error("Error opening file '%s': %s", file.Path, err.Error())
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, bytes.NewReader(file.Data))
	if err != nil {
		log.Error("Error writing file '%s': %s", file.Path, err.Error())
		return err
	}
	if err := f.Chown(file.UID, file.GID); err != nil {
		log.Error("Error chowning file '%s': %s", file.Path, err.Error())
		return err
	}

	log.Debug("Wrote %d bytes to '%s'", n, file.Path)

	return nil
}

func getCurrentUID() uint32 {
	current, err := user.Current()
	if err != nil {
		panic(err)
	}
	uid, err := strconv.ParseUint(current.Uid, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(uid)
}

func getCurrentGID() uint32 {
	current, err := user.Current()
	if err != nil {
		panic(err)
	}
	uid, err := strconv.ParseUint(current.Gid, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(uid)
}
