package main

import (
	"io"
	"os"
	"path"

	"github.com/ecnepsnai/logtic"
)

var log *logtic.Source

func main() {
	loadRegisterProperties()
	parseArgs()
	tryAutoRegister()
	mustLoadConfig()

	logtic.Log.FilePath = path.Join(config.LogPath, "otto_agent.log")
	if os.Getenv("OTTO_VERBOSE") == "" {
		logtic.Log.Stdout = io.Discard
	}

	logtic.Log.Open()
	log = logtic.Log.Connect("otto")

	mustLoadIdentity()
	setupLoopback()
	go startLoopbackRepeater()

	listen()
}
