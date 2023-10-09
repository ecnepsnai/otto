package main

import (
	"io"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/ecnepsnai/logtic"
)

var log *logtic.Source

func main() {
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-killSignal
		stop()
	}()

	printSignal := make(chan os.Signal, 1)
	signal.Notify(printSignal, syscall.SIGTRAP)
	go func() {
		<-printSignal
		runtime.Breakpoint()
	}()
	start()
}

func start() {
	if dir := os.Getenv("OTTO_DIR"); dir != "" {
		otto_DIR = dir
	}

	loadRegisterProperties()
	parseArgs()
	tryAutoRegister()
	mustLoadConfig()
	setupControl()
	go controlMain()

	logtic.Log.FilePath = path.Join(config.LogPath, "otto_agent.log")
	if os.Getenv("OTTO_VERBOSE") == "" {
		logtic.Log.Stdout = io.Discard
	}
	if os.Getenv("OTTO_DEBUG") != "" {
		logtic.Log.Level = logtic.LevelDebug
	}

	logtic.Log.Open()
	log = logtic.Log.Connect("otto")

	mustLoadIdentity()

	listen()
	stop()
}

func stop() {
	stopControl()
	os.Exit(1)
}
