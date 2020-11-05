package server

import (
	"path"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/security"
)

var neededTableVersion = 5

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
		if i == 5 {
			migrate5()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #5 Update user password
func migrate5() {
	log.Debug("Start migrate 5")

	if !FileExists(path.Join(Directories.Data, "user.db")) {
		return
	}

	type oldUser struct {
		Username     string `ds:"primary"`
		Email        string `ds:"unique"`
		Enabled      bool
		PasswordHash string
	}
	type newUser struct {
		Username     string `ds:"primary"`
		Email        string `ds:"unique"`
		Enabled      bool
		PasswordHash security.HashedPassword
	}

	result := ds.Migrate(ds.MigrateParams{
		TablePath: path.Join(Directories.Data, "user.db"),
		NewPath:   path.Join(Directories.Data, "user.db"),
		OldType:   oldUser{},
		NewType:   newUser{},
		MigrateObject: func(o interface{}) (interface{}, error) {
			oldUser := o.(oldUser)
			newHash := []byte(string(security.HashingAlgorithmBCrypt) + "$" + oldUser.PasswordHash)
			return User{
				Username:     oldUser.Username,
				Email:        oldUser.Email,
				Enabled:      oldUser.Enabled,
				PasswordHash: security.HashedPassword(newHash),
			}, nil
		},
	})
	if !result.Success {
		log.Fatal("Error migrating user table: %s", result.Error.Error())
	}
	log.Warn("Schedule store migration results: %+v", result)
}
