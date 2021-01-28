package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/ecnepsnai/otto/server/environ"
)

const configFileName = "otto_server.conf"

// OttoOptions describes options for the otto server
type OttoOptions struct {
	General        OptionsGeneral
	Authentication OptionsAuthentication
	Network        OptionsNetwork
	Register       OptionsRegister
	Security       OptionsSecurity
}

// OptionsGeneral describes the general options
type OptionsGeneral struct {
	ServerURL         string
	GlobalEnvironment []environ.Variable
}

// OptionsAuthentication describes the authentication options
type OptionsAuthentication struct {
	MaxAgeMinutes int
	SecureOnly    bool
}

// OptionsSecurity describes security options
type OptionsSecurity struct {
	IncludePSKEnv bool
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
	DefaultGroupID string
}

// Options the global options
var Options *OttoOptions
var optionsLock = sync.Mutex{}

// LoadOptions load Otto Server options
func LoadOptions() {
	defaults := OttoOptions{
		General: OptionsGeneral{
			ServerURL:         "http://" + bindAddress + "/",
			GlobalEnvironment: []environ.Variable{},
		},
		Authentication: OptionsAuthentication{
			MaxAgeMinutes: 60,
			SecureOnly:    false,
		},
		Network: OptionsNetwork{
			ForceIPVersion:     IPVersionOptionAuto,
			Timeout:            10,
			HeartbeatFrequency: 5,
		},
		Register: OptionsRegister{
			Enabled: false,
		},
		Security: OptionsSecurity{
			IncludePSKEnv: false,
		},
	}

	if !FileExists(path.Join(Directories.Data, configFileName)) {
		Options = &defaults
		Options.Save()
	} else {
		f, err := os.OpenFile(path.Join(Directories.Data, configFileName), os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal("Error opening config file: %s", err.Error())
		}
		defer f.Close()
		options := defaults
		if err := json.NewDecoder(f).Decode(&options); err != nil {
			log.Fatal("Error decoding options: %s", err.Error())
		}
		if err := options.Validate(); err != nil {
			log.Fatal("Invalid Otto Server Options: %s", err.Error())
		}
		Options = &options
	}
}

// Save save the options to disk. Will panic on any error. Returns true if the options did change
func (o *OttoOptions) Save() (string, bool) {
	optionsLock.Lock()
	defer optionsLock.Unlock()

	beforeHash := optionsFileHash()

	atomicPath := path.Join(Directories.Data, fmt.Sprintf(".%s_%s", configFileName, newPlainID()))
	realPath := path.Join(Directories.Data, configFileName)

	f, err := os.OpenFile(atomicPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Panic("Error opening config file: %s", err.Error())
	}
	if err := json.NewEncoder(f).Encode(o); err != nil {
		f.Close()
		log.Panic("Error encoding options: %s", err.Error())
	}
	f.Close()

	if err := os.Rename(atomicPath, realPath); err != nil {
		log.Panic("Error updating config file: %s", err.Error())
	}

	Options = o

	afterHash := optionsFileHash()
	return afterHash, beforeHash != afterHash
}

func optionsFileHash() string {
	configPath := path.Join(Directories.Data, configFileName)
	if !FileExists(configPath) {
		return ""
	}

	h, err := hashFile(configPath)
	if err != nil {
		log.Panic("Error hasing config file: %s", err.Error())
	}

	return h
}

// Validate returns an error if the options is not valid
func (o *OttoOptions) Validate() error {
	if o.General.ServerURL == "" {
		return fmt.Errorf("A server URL is required")
	}
	if err := environ.Validate(o.General.GlobalEnvironment); err != nil {
		return err
	}
	if !IsIPVersionOption(o.Network.ForceIPVersion) {
		return fmt.Errorf("Invalid value for IP version")
	}
	if o.Register.Enabled {
		if o.Register.PSK == "" {
			return fmt.Errorf("A register PSK is required if auto registration is enabled")
		}
	}
	return nil
}
