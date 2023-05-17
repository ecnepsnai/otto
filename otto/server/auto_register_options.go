package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
)

const autoregisterConfigFileName = "auto_register.conf"

// OptionsRegister describes register options
type OptionsRegister struct {
	Enabled        bool
	Key            string
	DefaultGroupID string
}

// AutoRegisterOptions the auto register options
var AutoRegisterOptions *OptionsRegister
var autoregisterOptionsLock = sync.Mutex{}

// LoadAutoRegisterOptions load auto register options
func LoadAutoRegisterOptions() {
	defaults := OptionsRegister{
		Enabled: false,
	}

	if !FileExists(path.Join(Directories.Data, configFileName)) {
		AutoRegisterOptions = &defaults
		AutoRegisterOptions.Save()
	} else {
		f, err := os.OpenFile(path.Join(Directories.Data, configFileName), os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal("Error opening config file: %s", err.Error())
		}
		defer f.Close()
		options := defaults
		if err := json.NewDecoder(f).Decode(&options); err != nil {
			log.Fatal("Error decoding autoregister options: %s", err.Error())
		}
		if err := options.Validate(); err != nil {
			log.Fatal("Invalid auto register options: %s", err.Error())
		}
		AutoRegisterOptions = &options
	}
}

// Save save the options to disk. Will panic on any error. Returns true if the options did change
func (o *OptionsRegister) Save() (string, bool) {
	autoregisterOptionsLock.Lock()
	defer autoregisterOptionsLock.Unlock()

	beforeHash := optionsFileHash()

	atomicPath := path.Join(Directories.Data, fmt.Sprintf(".%s_%s", autoregisterConfigFileName, newPlainID()))
	realPath := path.Join(Directories.Data, autoregisterConfigFileName)

	f, err := os.OpenFile(atomicPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Panic("Error opening autoregister config file: %s", err.Error())
	}
	if err := prettyJsonEncoder(f).Encode(o); err != nil {
		f.Close()
		log.Panic("Error encoding options: %s", err.Error())
	}
	f.Close()

	if err := os.Rename(atomicPath, realPath); err != nil {
		log.Panic("Error updating autoregister config file: %s", err.Error())
	}

	AutoRegisterOptions = o

	afterHash := optionsFileHash()
	return afterHash, beforeHash != afterHash
}

// Validate returns an error if the options is not valid
func (o *OptionsRegister) Validate() error {
	if o.Enabled {
		if o.Key == "" {
			return fmt.Errorf("a register key is required if auto registration is enabled")
		}
	}
	return nil
}
