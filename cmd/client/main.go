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
	logtic.Open()
	log = logtic.Connect("otto")

	l, err := net.Listen("tcp", "0.0.0.0:12444")
	if err != nil {
		panic(err)
	}
	log.Info("Otto client listening on 0.0.0.0:12444")
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go newRequest(c)
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
			fmt.Printf("Otto client %s (Runtime %s)\n", MainVersion, runtime.Version())
			os.Exit(0)
		}
		i++
	}
}

func newRequest(c net.Conn) {
	log.Info("New request from %s", c.RemoteAddr().String())
	defer c.Close()

	request, err := otto.ReadRequest(c, config.PSK)
	if err != nil {
		log.Error("Error reading request from '%s': %s", c.RemoteAddr().String(), err.Error())
		return
	}

	reply := otto.Reply{}
	switch request.Action {
	case otto.ActionPing, otto.ActionExit, otto.ActionReboot, otto.ActionShutdown:
		// No action
		break
	case otto.ActionReloadConfig:
		if err := loadConfig(); err != nil {
			reply.Error = err
		}
		break
	case otto.ActionRunScript:
		reply.ScriptResult = runScript(request.Script)
		break
	case otto.ActionUploadFile, otto.ActionUploadFileAndExit:
		if err := uploadFile(request.File); err != nil {
			reply.Error = err
		}
		break
	default:
		log.Error("Unknown action %d", request.Action)
		return
	}

	if err := otto.WriteReply(reply, config.PSK, c); err != nil {
		log.Error("Error writing reply to '%s': %s", c.RemoteAddr().String(), err.Error())
		return
	}

	if request.Action == otto.ActionUploadFileAndExit || request.Action == otto.ActionExit {
		c.Close()
		log.Warn("Exiting at request of '%s'", c.RemoteAddr().String())
		os.Exit(1)
	}

	if request.Action == otto.ActionReboot {
		c.Close()
		log.Warn("Rebooting at request of '%s'", c.RemoteAddr().String())
		exec.Command("/usr/sbin/reboot").Run()
	} else if request.Action == otto.ActionShutdown {
		c.Close()
		log.Warn("Shutting down at request of '%s'", c.RemoteAddr().String())
		exec.Command("/usr/sbin/halt").Run()
	}
}

func runScript(script otto.Script) otto.ScriptResult {
	tmp, err := ioutil.TempFile("", "otto")
	if err != nil {
		panic(err)
	}
	log.Debug("Writing script to %s", tmp.Name())
	if err := tmp.Chmod(0777); err != nil {
		panic(err)
	}
	defer os.Remove(tmp.Name())
	if _, err := io.CopyBuffer(tmp, bytes.NewBuffer(script.Data), nil); err != nil {
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

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Info("Running '%s' as UID %d GID %d", tmp.Name(), script.UID, script.GID)
	if err := cmd.Run(); err != nil {
		result.Success = false
		if exitError, ok := err.(*exec.ExitError); ok {
			result.Code = exitError.ExitCode()
			log.Error("Script exit code: %d", result.Code)
		} else {
			log.Error("Error running script: %s", err.Error())
			result.ExecError = err.Error()
		}
	} else {
		log.Info("Script exit OK")
		result.Success = true
	}

	result.Stderr = string(stderr.Bytes())
	result.Stdout = string(stdout.Bytes())
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
