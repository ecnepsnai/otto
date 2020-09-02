package server

import (
	"path"

	"github.com/ecnepsnai/logtic"
)

// Start the app
func Start() {
	preBootstrapArgs()
	startup()
	postBootstrapArgs()
	RouterSetup()
}

// Stop stop the API service gracefully
func Stop() {
	shutdown()
}

var log *logtic.Source

// CommonSetup common setup methods
func CommonSetup() {
	fsSetup()
	initLogtic(isVerbose())
	GobSetup()
	StateSetup()
	migrateIfNeeded()
	LoadOptions()
}

func initLogtic(verbose bool) {
	logtic.Log.Level = logtic.LevelInfo
	if verbose {
		logtic.Log.Level = logtic.LevelDebug
	}
	logtic.Log.FilePath = path.Join(Directories.Logs, "otto.log")
	if err := logtic.Open(); err != nil {
		panic(err)
	}
	log = logtic.Connect("otto")
}

func startup() {
	CommonSetup()
	DataStoreSetup()
	WarmCache()
	ScheduleSetup()
	checkFirstRun()
	go StartHeartbeatMonitor()
}

func shutdown() {
	State.Close()
	DataStoreTeardown()
	logtic.Close()
}
