package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"sync"

	"github.com/ecnepsnai/otto/server/environ"
)

// OttoOptions describes options for the otto server
type OttoOptions struct {
	General  OptionsGeneral
	Network  OptionsNetwork
	Register OptionsRegister
	Security OptionsSecurity
}

// OptionsGeneral describes the general options
type OptionsGeneral struct {
	ServerURL         string
	GlobalEnvironment []environ.Variable
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
	Rules          []RegisterRule
	DefaultGroupID string
}

// RegisterRule describes a register rule
type RegisterRule struct {
	Property string
	Pattern  string
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
		Security: OptionsSecurity{
			IncludePSKEnv: false,
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
		if err := options.Validate(); err != nil {
			log.Fatal("Invalid Otto Server Options: %s", err.Error())
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
		for _, rule := range o.Register.Rules {
			if rule.GroupID == "" {
				return fmt.Errorf("Invalid group ID on registration rule")
			}
			if !IsRegisterRuleProperty(rule.Property) {
				return fmt.Errorf("Invalid registration rule property")
			}
			if rule.Pattern == "" {
				return fmt.Errorf("Missing registration rule pattern")
			}
			if _, err := regexp.Compile(rule.Pattern); err != nil {
				return fmt.Errorf("Invalid regex pattern on registration rule")
			}
		}
	}
	return nil
}
