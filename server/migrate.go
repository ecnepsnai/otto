package server

import (
	"path"
	"sync"
	"time"

	"github.com/ecnepsnai/ds"
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
func migrate10() {
	log.Debug("Start migrate 10")

	wg := sync.WaitGroup{}
	wg.Add(2)

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

	wg.Wait()
	StoreTeardown()

	if !allSuccess {
		log.Fatal("One or more migrations failed, see above log lines for more information")
	}
}
