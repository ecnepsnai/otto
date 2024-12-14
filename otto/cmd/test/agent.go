package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func registerHost() error {
	return logTask("Register host", func() error {
		if err := os.MkdirAll(path.Join(dataDir, "agent"), 7644); err != nil {
			return err
		}

		cmd := exec.Command(path.Join(testDir, "agent"))
		cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
		cmd.Env = append(cmd.Env, "OTTO_DEBUG=1")
		cmd.Env = append(cmd.Env, "OTTO_VERBOSE=1")
		cmd.Env = append(cmd.Env, "REGISTER_HOST=http://127.0.0.1:8080")
		cmd.Env = append(cmd.Env, fmt.Sprintf("REGISTER_KEY=%s", registerKey))
		cmd.Dir = path.Join(dataDir, "agent")

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf(string(output))
		}
		return nil
	})
}

func startAgent(stop chan uint8) error {
	if err := os.MkdirAll(path.Join(dataDir, "agent"), 7644); err != nil {
		return err
	}

	out, err := os.OpenFile(path.Join(dataDir, "agent", "stdout.txt"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command(path.Join(testDir, "agent"))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	cmd.Env = append(cmd.Env, "OTTO_DEBUG=1")
	cmd.Env = append(cmd.Env, "OTTO_VERBOSE=1")
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Dir = path.Join(dataDir, "agent")
	if err := cmd.Start(); err != nil {
		return err
	}

	process := cmd.Process
	_ = <-stop
	process.Signal(syscall.SIGINT)
	fmt.Println("Stopping agent")
	out.Sync()
	out.Close()
	return nil
}
