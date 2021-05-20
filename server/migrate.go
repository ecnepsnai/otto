package server

import (
	"encoding/json"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
	"github.com/ecnepsnai/secutil"
)

var neededTableVersion = 10

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
		if i == 10 {
			migrate10()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #10
// - move user password & api key hash into dedicated store
// - update script attachment owner
// - update script run as type
// - rename register PSK as Key
func migrate10() {
	log.Debug("Start migrate 10")

	wg := sync.WaitGroup{}
	wg.Add(4)

	allSuccess := true

	StoreSetup()

	go func() {
		defer wg.Done()

		if !FileExists(path.Join(Directories.Data, "user.db")) {
			return
		}

		type oldUserType struct {
			Username           string `ds:"primary" max:"32" min:"1"`
			Email              string `ds:"unique" max:"128" min:"1"`
			APIKey             secutil.HashedPassword
			PasswordHash       secutil.HashedPassword
			CanLogIn           bool
			MustChangePassword bool
		}
		type newUserType struct {
			Username           string `ds:"primary" max:"32" min:"1"`
			Email              string `ds:"unique" max:"128" min:"1"`
			CanLogIn           bool
			MustChangePassword bool
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
					CanLogIn:           oldUser.CanLogIn,
					MustChangePassword: oldUser.MustChangePassword,
				}
				ShadowStore.Set(newUser.Username, oldUser.PasswordHash)
				if len(oldUser.APIKey) > 0 {
					ShadowStore.Set("api_"+newUser.Username, oldUser.APIKey)
				}
				return newUser, nil
			},
		})

		if results.Error != nil {
			log.Error("Error migrating user database: %s", results.Error.Error())
			allSuccess = false
		}
	}()

	go func() {
		defer wg.Done()

		if !FileExists(path.Join(Directories.Data, "attachment.db")) {
			return
		}

		type oldAttachmentType struct {
			ID       string `ds:"primary"`
			Path     string
			Name     string
			MimeType string
			UID      int
			GID      int
			Created  time.Time
			Modified time.Time
			Mode     uint32
			Size     uint64
		}

		type newAttachmentType struct {
			ID       string `ds:"primary"`
			Path     string
			Name     string
			MimeType string
			Owner    RunAs
			Created  time.Time
			Modified time.Time
			Mode     uint32
			Size     uint64
		}

		results := ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "attachment.db"),
			NewPath:   path.Join(Directories.Data, "attachment.db"),
			OldType:   oldAttachmentType{},
			NewType:   newAttachmentType{},
			MigrateObject: func(old interface{}) (interface{}, error) {
				oldAttachment, ok := old.(oldAttachmentType)
				if !ok {
					panic("Invalid type")
				}
				newAttachment := newAttachmentType{
					ID:       oldAttachment.ID,
					Path:     oldAttachment.Path,
					Name:     oldAttachment.Name,
					MimeType: oldAttachment.MimeType,
					Owner: RunAs{
						Inherit: false,
						UID:     uint32(oldAttachment.UID),
						GID:     uint32(oldAttachment.GID),
					},
					Created:  oldAttachment.Created,
					Modified: oldAttachment.Modified,
					Mode:     oldAttachment.Mode,
					Size:     oldAttachment.Size,
				}
				return newAttachment, nil
			},
		})

		if results.Error != nil {
			log.Error("Error migrating attachment database: %s", results.Error.Error())
			allSuccess = false
		}
	}()

	go func() {
		defer wg.Done()

		if !FileExists(path.Join(Directories.Data, "script.db")) {
			return
		}

		type ScriptRunAs struct {
			Inherit bool
			UID     uint32
			GID     uint32
		}

		type oldScriptType struct {
			ID               string `ds:"primary"`
			Name             string `ds:"unique" min:"1" max:"140"`
			Enabled          bool   `ds:"index"`
			Executable       string `min:"1"`
			Script           string `min:"1"`
			Environment      []environ.Variable
			RunAs            ScriptRunAs
			WorkingDirectory string
			AfterExecution   string
			AttachmentIDs    []string
		}

		type newScriptType struct {
			ID               string `ds:"primary"`
			Name             string `ds:"unique" min:"1" max:"140"`
			Enabled          bool   `ds:"index"`
			Executable       string `min:"1"`
			Script           string `min:"1"`
			Environment      []environ.Variable
			RunAs            RunAs
			WorkingDirectory string
			AfterExecution   string
			AttachmentIDs    []string
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
					RunAs: RunAs{
						Inherit: oldScript.RunAs.Inherit,
						UID:     oldScript.RunAs.UID,
						GID:     oldScript.RunAs.GID,
					},
					WorkingDirectory: oldScript.WorkingDirectory,
					AfterExecution:   oldScript.AfterExecution,
					AttachmentIDs:    oldScript.AttachmentIDs,
				}
				return newScript, nil
			},
		})

		if results.Error != nil {
			log.Error("Error migrating script database: %s", results.Error.Error())
			allSuccess = false
		}
	}()

	go func() {
		defer wg.Done()

		if !FileExists(path.Join(Directories.Data, configFileName)) {
			return
		}

		if FileExists(path.Join(Directories.Data, configFileName) + "_backup") {
			return
		}

		f, err := os.OpenFile(path.Join(Directories.Data, configFileName), os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Error("Error opening config file: %s", err.Error())
			allSuccess = false
			return
		}

		optionsMap := map[string]interface{}{}
		if err := json.NewDecoder(f).Decode(&optionsMap); err != nil {
			log.Error("Error decoding options file: %s", err.Error())
			allSuccess = false
			f.Close()
			return
		}

		register, present := optionsMap["Register"].(map[string]interface{})
		if !present {
			f.Close()
			return
		}

		psk, present := register["PSK"]
		if !present {
			f.Close()
			return
		}

		delete(register, "PSK")
		register["Key"] = psk
		optionsMap["Register"] = register

		if err := os.Rename(path.Join(Directories.Data, configFileName), path.Join(Directories.Data, configFileName)+"_backup"); err != nil {
			log.Error("Error making backup of config file: %s", err.Error())
			allSuccess = false
			return
		}

		f, err = os.OpenFile(path.Join(Directories.Data, configFileName), os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Error("Error opening config file: %s", err.Error())
			allSuccess = false
			return
		}

		if err := json.NewEncoder(f).Encode(&optionsMap); err != nil {
			log.Error("Error writing new options file: %s", err.Error())
			allSuccess = false
			f.Close()
		}
	}()

	wg.Wait()
	StoreTeardown()

	if !allSuccess {
		log.Fatal("One or more migrations failed, see above log lines for more information")
	}
}
