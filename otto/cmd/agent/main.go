package main

import (
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/ecnepsnai/logtic"
)

var log *logtic.Source

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		stop()
	}()
	start()
}

func start() {
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
	setupLoopback()
	go startLoopbackRepeater()

	listen()
	stop()
}

func stop() {
	stopControl()
	os.Exit(1)
}
