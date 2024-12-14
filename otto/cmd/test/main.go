package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
)

var dockerCmd = "podman"

var testDir = ""
var dataDir = ""
var cmdDir = ""
var ottoDir = ""
var serverApiKey = ""
var registerKey = "REGISTER_KEY"

func main() {
	if dockerCmdEnv := os.Getenv("DOCKER_CMD"); dockerCmdEnv != "" {
		dockerCmd = dockerCmdEnv
	}

	resolveDirs()

	if !runTest() {
		os.Exit(1)
	}
}

func runTest() bool {
	if compileAgent() != nil {
		return false
	}
	if compileServer() != nil {
		return false
	}

	serverStop := make(chan uint8)
	go startServer(serverStop)
	defer func() {
		serverStop <- uint8(1)
		time.Sleep(100 * time.Millisecond)
	}()

	// Wait for the server to start
	time.Sleep(500 * time.Millisecond)

	// Get the API key
	apiKey, err := getApiKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating API key: %s\n", err.Error())
		return false
	}
	serverApiKey = *apiKey

	// Configure automatic registration
	if err := enableAutoRegister(); err != nil {
		fmt.Fprintf(os.Stderr, "Error configuring automatic host registration: %s\n", err.Error())
		return false
	}

	// Create Script
	if err := createScript(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating test script: %s\n", err.Error())
		return false
	}

	// Register Host
	if err := registerHost(); err != nil {
		fmt.Fprintf(os.Stderr, "Error registering host: %s\n", err.Error())
		return false
	}

	agentStop := make(chan uint8)
	go startAgent(agentStop)
	defer func() {
		agentStop <- uint8(1)
		time.Sleep(100 * time.Millisecond)
	}()

	// Wait for the agent to start
	time.Sleep(500 * time.Millisecond)

	// Run a script
	if err := runScript(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running script: %s\n", err.Error())
		return false
	}

	// Rotate the client identity
	if err := rotateIdentity(); err != nil {
		fmt.Fprintf(os.Stderr, "Error rotating client identity: %s\n", err.Error())
		return false
	}

	// Wait for the agent to start
	time.Sleep(500 * time.Millisecond)

	// Run a script
	if err := runScript(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running script: %s\n", err.Error())
		return false
	}

	return true
}

func resolveDirs() {
	var err error
	testDir, err = filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	dataDir = path.Join(testDir, "data")
	if err := os.MkdirAll(dataDir, 7644); err != nil {
		panic(err)
	}
	cmdDir, err = filepath.Abs("../")
	if err != nil {
		panic(err)
	}
	ottoDir, err = filepath.Abs("../../")
	if err != nil {
		panic(err)
	}
}

func logTask(label string, task func() error) error {
	fmt.Printf("** %s...", label)
	err := task()
	if err != nil {
		fmt.Printf(" [\x1b[31mFAIL\x1b[0m]\n")
		return err
	}
	fmt.Printf(" [\x1b[32mDONE\x1b[0m]\n")
	return nil
}
