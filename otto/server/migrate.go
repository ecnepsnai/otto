package server

import (
	"encoding/json"
	"os"
	"path"

	"github.com/ecnepsnai/ds"
)

var neededTableVersion = 13

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
		i++

		results := ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "user.db"),
			NewPath:   path.Join(Directories.Data, "user.db"),
			OldType:   User{},
			NewType:   User{},
			MigrateObject: func(old interface{}) (interface{}, error) {
				user, ok := old.(User)
				if !ok {
					panic("Invalid type")
				}
				user.Permissions = UserPermissionsMax()
				return user, nil
			},
		})

		if results.Error != nil {
			log.Fatal("Error migrating user database: %s", results.Error.Error())
		}

		results = ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "script.db"),
			NewPath:   path.Join(Directories.Data, "script.db"),
			OldType:   Script{},
			NewType:   Script{},
			MigrateObject: func(old interface{}) (interface{}, error) {
				script, ok := old.(Script)
				if !ok {
					panic("Invalid type")
				}
				script.RunLevel = ScriptRunLevelReadWrite
				return script, nil
			},
		})

		if results.Error != nil {
			log.Fatal("Error migrating script database: %s", results.Error.Error())
		}

		if FileExists(path.Join(Directories.Data, configFileName)) {
			type oldOttoOptions struct {
				General        OptionsGeneral
				Authentication OptionsAuthentication
				Network        OptionsNetwork
				Register       OptionsRegister
				Security       OptionsSecurity
			}
			f, err := os.OpenFile(path.Join(Directories.Data, configFileName), os.O_RDONLY, os.ModePerm)
			if err != nil {
				log.Fatal("Unable to open Otto config file: %s", err.Error())
			}
			oldOptions := oldOttoOptions{}
			if err := json.NewDecoder(f).Decode(&oldOptions); err != nil {
				log.Fatal("Unable to open Otto config file: %s", err.Error())
			}
			oldOptions.Register.Save()
			newOptions := OttoOptions{
				General:        oldOptions.General,
				Authentication: oldOptions.Authentication,
				Network:        oldOptions.Network,
				Security:       oldOptions.Security,
			}
			newOptions.Save()
		}
	}

	State.SetTableVersion(i)
}
