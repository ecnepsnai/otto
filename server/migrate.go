package server

import (
	"path"

	"github.com/ecnepsnai/ds"
)

var neededTableVersion = 6

func migrateIfNeeded() {
	currentVersion := State.GetTableVersion()

	if currentVersion == 0 {
		State.SetTableVersion(neededTableVersion + 1)
		log.Debug("Setting default table version to %d", neededTableVersion+1)
		return
	}

	if neededTableVersion-currentVersion > 1 {
		log.Fatal("Refusing to migrate datastore that is too old - follow the supported upgrade path and don't skip versions. Table version %d, required version %d", currentVersion, neededTableVersion)
	}

	i := currentVersion
	for i <= neededTableVersion {
		if i == 6 {
			migrate6()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #6 Update schedule report
func migrate6() {
	log.Debug("Start migrate 6")

	if !FileExists(path.Join(Directories.Data, "schedulereport.db")) {
		return
	}

	type scheduleReport struct {
		ID         string `ds:"primary"`
		ScheduleID string `ds:"index"`
		HostIDs    []string
		Time       ScheduleReportTime
		Result     int
		HostResult map[string]int
	}

	result := ds.Migrate(ds.MigrateParams{
		TablePath: path.Join(Directories.Data, "schedulereport.db"),
		NewPath:   path.Join(Directories.Data, "schedulereport.db"),
		OldType:   scheduleReport{},
		NewType:   scheduleReport{},
		MigrateObject: func(o interface{}) (interface{}, error) {
			report := o.(scheduleReport)
			report.HostResult = map[string]int{}
			for _, hostID := range report.HostIDs {
				report.HostResult[hostID] = -1
			}
			return report, nil
		},
	})
	if !result.Success {
		log.Fatal("Error migrating schedulereport table: %s", result.Error.Error())
	}
	log.Warn("Schedule store migration results: %+v", result)
}
