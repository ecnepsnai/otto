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
	LoadOptions()
}

// StatsPopulate populate the stats object
func StatsPopulate() {

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
	StatsPopulate()
	ScheduleSetup()
	checkFirstRun()
}

func shutdown() {
	DataStoreTeardown()
	logtic.Close()
}
