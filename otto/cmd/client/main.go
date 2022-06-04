package main

import (
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
	mustLoadIdentity()

	logtic.Log.FilePath = path.Join(config.LogPath, "otto_client.log")
	logtic.Log.Level = logtic.LevelWarn
	if os.Getenv("OTTO_VERBOSE") != "" {
		logtic.Log.Level = logtic.LevelDebug
	}

	logtic.Log.Open()
	log = logtic.Log.Connect("otto")

	setupLoopback()
	go startLoopbackRepeater()

	listen()
}
