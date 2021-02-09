package server

import (
	"path"
	"sync"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
	"github.com/ecnepsnai/secutil"
)

var neededTableVersion = 8

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
		if i == 8 {
			migrate8()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #8 update script run as
func migrate8() {
	log.Debug("Start migrate 8")

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		type oldScriptType struct {
			ID               string `ds:"primary"`
			Name             string `ds:"unique"`
			Enabled          bool   `ds:"index"`
			Executable       string
			Script           string
			Environment      []environ.Variable
			UID              uint32
			GID              uint32
			WorkingDirectory string
			AfterExecution   string
			AttachmentIDs    []string
		}

		type newScriptType struct {
			ID               string `ds:"primary"`
			Name             string `ds:"unique"`
			Enabled          bool   `ds:"index"`
			Executable       string
			Script           string
			Environment      []environ.Variable
			RunAs            ScriptRunAs
			WorkingDirectory string
			AfterExecution   string
			AttachmentIDs    []string
		}

		if !FileExists(path.Join(Directories.Data, "script.db")) {
			return
		}

		results := ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "script.db"),
			NewPath:   path.Join(Directories.Data, "script.db"),
			OldType:   oldScriptType{},
			NewType:   newScriptType{},
			MigrateObject: func(old interface{}) (interface{}, error) {
				oldScript, ok := old.(oldScriptType)
				if !ok {
					panic("Invalid type")
				}
				newScript := newScriptType{
					ID:          oldScript.ID,
					Name:        oldScript.Name,
					Enabled:     oldScript.Enabled,
					Executable:  oldScript.Executable,
					Script:      oldScript.Script,
					Environment: oldScript.Environment,
					RunAs: ScriptRunAs{
						Inherit: false,
						UID:     oldScript.UID,
						GID:     oldScript.GID,
					},
					WorkingDirectory: oldScript.WorkingDirectory,
					AfterExecution:   oldScript.AfterExecution,
					AttachmentIDs:    oldScript.AttachmentIDs,
				}
				return newScript, nil
			},
		})

		if results.Error != nil {
			log.Fatal("Error migrating script database: %s", results.Error.Error())
		}

		log.Debug("Migrated script: %+v", results)
	}()

	go func() {
		defer wg.Done()

		type oldUserType struct {
			Username     string `ds:"primary"`
			Email        string `ds:"unique"`
			PasswordHash secutil.HashedPassword
			Enabled      bool
		}
		type newUserType struct {
			Username           string `ds:"primary"`
			Email              string `ds:"unique"`
			PasswordHash       secutil.HashedPassword
			CanLogIn           bool
			MustChangePassword bool
		}

		if !FileExists(path.Join(Directories.Data, "user.db")) {
			return
		}

		results := ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "user.db"),
			NewPath:   path.Join(Directories.Data, "user.db"),
			OldType:   oldUserType{},
			NewType:   newUserType{},
			MigrateObject: func(old interface{}) (interface{}, error) {
				oldUser, ok := old.(oldUserType)
				if !ok {
					panic("Invalid type")
				}
				newUser := newUserType{
					Username:           oldUser.Username,
					Email:              oldUser.Email,
					PasswordHash:       oldUser.PasswordHash,
					CanLogIn:           oldUser.Enabled,
					MustChangePassword: false,
				}

				if oldUser.PasswordHash.Compare([]byte("")) {
					newUser.MustChangePassword = true
				}

				return newUser, nil
			},
		})

		if results.Error != nil {
			log.Fatal("Error migrating user database: %s", results.Error.Error())
		}

		log.Debug("Migrated user: %+v", results)
	}()

	wg.Wait()
}
