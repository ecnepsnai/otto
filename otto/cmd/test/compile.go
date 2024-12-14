package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func compileServer() error {
	return logTask("Compiling server", func() error {
		cmd := exec.Command("go", "build")
		cmd.Dir = path.Join(cmdDir, "server")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf(string(output))
		}
		os.Rename(path.Join(cmd.Dir, "server"), path.Join(testDir, "server"))
		return nil
	})
}

func compileAgent() error {
	return logTask("Compiling agent", func() error {
		cmd := exec.Command("go", "build")
		cmd.Dir = path.Join(cmdDir, "agent")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf(string(output))
		}
		os.Rename(path.Join(cmd.Dir, "agent"), path.Join(testDir, "agent"))
		return nil
	})
}
