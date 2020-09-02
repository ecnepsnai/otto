package server

import (
	"encoding/json"
	"os"
	"path"
	"sync"

	"github.com/ecnepsnai/otto/server/environ"
)

// OttoOptions describes options for the otto server
type OttoOptions struct {
	General  OptionsGeneral
	Network  OptionsNetwork
	Register OptionsRegister
}

// OptionsGeneral describes the general options
type OptionsGeneral struct {
	ServerURL         string
	GlobalEnvironment []environ.Variable
}

// OptionsNetwork describes network options for connecting to otto clients
type OptionsNetwork struct {
	ForceIPVersion     string
	Timeout            int64
	HeartbeatFrequency int64
}

// OptionsRegister describes register options
type OptionsRegister struct {
	Enabled        bool
	PSK            string
	Rules          []RegisterRule
	DefaultGroupID string
}

// RegisterRule describes a register rule
type RegisterRule struct {
	Uname    string
	Hostname string
	GroupID  string
}

// Options the global options
var Options *OttoOptions
var optionsLock = sync.Mutex{}

// LoadOptions load E6 options
func LoadOptions() {
	defaults := OttoOptions{
		General: OptionsGeneral{
			ServerURL:         "http://" + bindAddress + "/",
			GlobalEnvironment: []environ.Variable{},
		},
		Network: OptionsNetwork{
			ForceIPVersion:     IPVersionOptionAuto,
			Timeout:            10,
			HeartbeatFrequency: 5,
		},
		Register: OptionsRegister{
			Enabled: false,
			Rules:   []RegisterRule{},
		},
	}

	if !FileExists(path.Join(Directories.Data, "otto_server.conf")) {
		Options = &defaults
		if err := Options.Save(); err != nil {
			log.Fatal("Error setting default options: %s", err.Error())
		}
	} else {
		f, err := os.OpenFile(path.Join(Directories.Data, "otto_server.conf"), os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal("Error opening config file: %s", err.Error())
		}
		defer f.Close()
		options := defaults
		if err := json.NewDecoder(f).Decode(&options); err != nil {
			log.Fatal("Error decoding options: %s", err.Error())
		}
		Options = &options
	}
}

// Save save the options to disk
func (o *OttoOptions) Save() error {
	optionsLock.Lock()
	defer optionsLock.Unlock()

	f, err := os.OpenFile(path.Join(Directories.Data, "otto_server.conf"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error("Error opening config file: %s", err.Error())
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(o); err != nil {
		log.Error("Error encoding options: %s", err.Error())
		return err
	}

	Options = o

	return nil
}
