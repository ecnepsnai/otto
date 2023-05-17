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
	gobSetup()
	stateSetup()
	migrateIfNeeded()
	LoadOptions()
	LoadAutoRegisterOptions()
	statsSetup()
}

func initLogtic(verbose bool) {
	logtic.Log.Level = logtic.LevelInfo
	if verbose {
		logtic.Log.Level = logtic.LevelDebug
	}
	logtic.Log.FilePath = path.Join(Directories.Logs, "otto.log")
	if err := logtic.Log.Open(); err != nil {
		panic(err)
	}
	log = logtic.Log.Connect("otto")
}

func startup() {
	CommonSetup()
	storeSetup()
	dataStoreSetup()
	AttachmentStore.Cleanup()
	CacheSetup()
	CronSetup()
	checkFirstRun()
	go StartHeartbeatMonitor()
}

func shutdown() {
	State.Close()
	dataStoreTeardown()
	storeTeardown()
	logtic.Log.Close()
}
