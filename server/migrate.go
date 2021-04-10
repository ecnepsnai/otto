package server

import (
	"path"

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

// #10 move user password hash into dedicated store
func migrate10() {
	log.Debug("Start migrate 10")

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

	if !FileExists(path.Join(Directories.Data, "user.db")) {
		return
	}

	StoreSetup()

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
		log.Fatal("Error migrating user database: %s", results.Error.Error())
	}

	StoreTeardown()

	log.Debug("Migrated user: %+v", results)
}
