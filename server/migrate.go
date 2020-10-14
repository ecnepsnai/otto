package server

import (
	"path"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/security"
)

var neededTableVersion = 4

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
		if i == 4 {
			migrate4()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #4 Add schedule names
func migrate4() {
	log.Debug("Start migrate 4")

	if !FileExists(path.Join(Directories.Data, "schedule.db")) {
		return
	}

	result := ds.Migrate(ds.MigrateParams{
		TablePath: path.Join(Directories.Data, "schedule.db"),
		NewPath:   path.Join(Directories.Data, "schedule.db"),
		OldType:   Schedule{},
		NewType:   Schedule{},
		MigrateObject: func(o interface{}) (interface{}, error) {
			schedule := o.(Schedule)
			var name = security.RandomString(5)
			script, _ := ScriptStore.ScriptWithID(schedule.ScriptID)
			if script != nil {
				name = script.Name
			}
			schedule.Name = name
			return schedule, nil
		},
	})
	if !result.Success {
		log.Fatal("Error migrating schedule table: %s", result.Error.Error())
	}
	log.Warn("Schedule store migration results: %+v", result)
}
