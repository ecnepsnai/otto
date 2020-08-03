package server

import (
	"encoding/json"
	"os"
	"path"
)

var neededTableVersion = 11

func migrateIfNeeded() {
	currentVersion := State.GetTableVersion()

	i := currentVersion
	for i <= neededTableVersion {
		if i == 1 {
			migrate1()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #1 organize options
func migrate1() {
	type lastOttoOptions struct {
		ServerURL         string
		GlobalEnvironment map[string]string
		Network           OptionsNetwork
		Register          OptionsRegister
	}

	if !FileExists(path.Join(Directories.Data, "otto_server.conf")) {
		return
	}

	f, err := os.OpenFile(path.Join(Directories.Data, "otto_server.conf"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal("Error opening config file: %s", err.Error())
	}
	defer f.Close()
	oldOptions := lastOttoOptions{}
	if err := json.NewDecoder(f).Decode(&oldOptions); err != nil {
		log.Fatal("Error decoding options: %s", err.Error())
	}

	newOptions := OttoOptions{
		General: OptionsGeneral{
			ServerURL:         oldOptions.ServerURL,
			GlobalEnvironment: oldOptions.GlobalEnvironment,
		},
		Network:  oldOptions.Network,
		Register: oldOptions.Register,
	}

	if err := newOptions.Save(); err != nil {
		log.Fatal("Error migrating options: %s", err.Error())
	}
}
