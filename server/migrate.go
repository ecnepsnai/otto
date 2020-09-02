package server

import (
	"encoding/json"
	"os"
	"path"
	"sync"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
)

var neededTableVersion = 4

func migrateIfNeeded() {
	currentVersion := State.GetTableVersion()

	i := currentVersion
	for i <= neededTableVersion {
		if i == 1 {
			migrate1()
		} else if i == 2 {
			migrate2()
		} else if i == 3 {
			migrate3()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #1 organize options
func migrate1() {
	log.Debug("Start migrate 1")
	type oldOptionsNetwork struct {
		ForceIPVersion string
		Timeout        int64
	}

	type oldRegisterRule struct {
		Uname    string
		Hostname string
		GroupID  string
	}

	type oldOptionsRegister struct {
		Enabled        bool
		PSK            string
		Rules          []oldRegisterRule
		DefaultGroupID string
	}

	type oldOttoOptions struct {
		ServerURL         string
		GlobalEnvironment map[string]string
		Network           oldOptionsNetwork
		Register          oldOptionsRegister
	}

	type newOptionsGeneral struct {
		ServerURL         string
		GlobalEnvironment map[string]string
	}

	type newOptionsNetwork struct {
		ForceIPVersion string
		Timeout        int64
	}

	type newRegisterRule struct {
		Uname    string
		Hostname string
		GroupID  string
	}

	type newOptionsRegister struct {
		Enabled        bool
		PSK            string
		Rules          []newRegisterRule
		DefaultGroupID string
	}

	type newOttoOptions struct {
		General  newOptionsGeneral
		Network  newOptionsNetwork
		Register newOptionsRegister
	}

	if !FileExists(path.Join(Directories.Data, "otto_server.conf")) {
		return
	}

	f, err := os.OpenFile(path.Join(Directories.Data, "otto_server.conf"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal("Error opening config file: %s", err.Error())
	}
	oldOptions := oldOttoOptions{}
	if err := json.NewDecoder(f).Decode(&oldOptions); err != nil {
		log.Fatal("Error decoding options: %s", err.Error())
	}
	f.Close()

	newOptions := newOttoOptions{
		General: newOptionsGeneral{
			ServerURL:         oldOptions.ServerURL,
			GlobalEnvironment: oldOptions.GlobalEnvironment,
		},
		Network: newOptionsNetwork{
			ForceIPVersion: oldOptions.Network.ForceIPVersion,
			Timeout:        oldOptions.Network.Timeout,
		},
		Register: newOptionsRegister{
			Enabled: oldOptions.Register.Enabled,
			Rules:   make([]newRegisterRule, len(oldOptions.Register.Rules)),
		},
	}

	for i, rule := range oldOptions.Register.Rules {
		newOptions.Register.Rules[i] = newRegisterRule{
			Uname:    rule.Uname,
			Hostname: rule.Hostname,
			GroupID:  rule.GroupID,
		}
	}

	f, err = os.OpenFile(path.Join(Directories.Data, "otto_server.conf"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal("Error migrating options: %s", err.Error())
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(newOptions); err != nil {
		log.Fatal("Error migrating options: %s", err.Error())
	}
}

// #2 update global environment variables
func migrate2() {
	log.Debug("Start migrate 2")
	type oldOptionsGeneral struct {
		ServerURL         string
		GlobalEnvironment map[string]string
	}

	type oldOptionsNetwork struct {
		ForceIPVersion     string
		Timeout            int64
		HeartbeatFrequency int64
	}

	type oldRegisterRule struct {
		Uname    string
		Hostname string
		GroupID  string
	}

	type oldOptionsRegister struct {
		Enabled        bool
		PSK            string
		Rules          []oldRegisterRule
		DefaultGroupID string
	}

	type oldOttoOptions struct {
		General  oldOptionsGeneral
		Network  oldOptionsNetwork
		Register oldOptionsRegister
	}

	type newOptionsGeneral struct {
		ServerURL         string
		GlobalEnvironment []environ.Variable
	}

	type newOptionsNetwork struct {
		ForceIPVersion     string
		Timeout            int64
		HeartbeatFrequency int64
	}

	type newRegisterRule struct {
		Uname    string
		Hostname string
		GroupID  string
	}

	type newOptionsRegister struct {
		Enabled        bool
		PSK            string
		Rules          []newRegisterRule
		DefaultGroupID string
	}

	type newOttoOptions struct {
		General  newOptionsGeneral
		Network  newOptionsNetwork
		Register newOptionsRegister
	}

	if !FileExists(path.Join(Directories.Data, "otto_server.conf")) {
		return
	}

	f, err := os.OpenFile(path.Join(Directories.Data, "otto_server.conf"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal("Error opening config file: %s", err.Error())
	}
	oldOptions := oldOttoOptions{}
	if err := json.NewDecoder(f).Decode(&oldOptions); err != nil {
		log.Fatal("Error decoding options: %s", err.Error())
	}
	f.Close()

	newOptions := newOttoOptions{
		General: newOptionsGeneral{
			ServerURL:         oldOptions.General.ServerURL,
			GlobalEnvironment: environ.FromMap(oldOptions.General.GlobalEnvironment),
		},
		Network: newOptionsNetwork{
			ForceIPVersion:     oldOptions.Network.ForceIPVersion,
			Timeout:            oldOptions.Network.Timeout,
			HeartbeatFrequency: oldOptions.Network.HeartbeatFrequency,
		},
		Register: newOptionsRegister{
			Enabled: oldOptions.Register.Enabled,
			Rules:   make([]newRegisterRule, len(oldOptions.Register.Rules)),
		},
	}

	for i, rule := range oldOptions.Register.Rules {
		newOptions.Register.Rules[i] = newRegisterRule{
			Uname:    rule.Uname,
			Hostname: rule.Hostname,
			GroupID:  rule.GroupID,
		}
	}

	f, err = os.OpenFile(path.Join(Directories.Data, "otto_server.conf"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal("Error migrating options: %s", err.Error())
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(newOptions); err != nil {
		log.Fatal("Error migrating options: %s", err.Error())
	}
}

// #3 update host, group, script environment variables
func migrate3() {
	log.Debug("Start migrate 3")
	type oldHost struct {
		ID          string `ds:"primary"`
		Name        string `ds:"unique"`
		Address     string `ds:"unique"`
		Port        uint32
		PSK         string
		Enabled     bool `ds:"index"`
		GroupIDs    []string
		Environment map[string]string
	}

	type oldScript struct {
		ID               string `ds:"primary"`
		Name             string `ds:"unique"`
		Enabled          bool   `ds:"index"`
		Executable       string
		Script           string
		Environment      map[string]string
		UID              uint32
		GID              uint32
		WorkingDirectory string
		AfterExecution   string
	}

	type oldGroup struct {
		ID          string `ds:"primary"`
		Name        string `ds:"unique"`
		ScriptIDs   []string
		Environment map[string]string
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		if !FileExists(path.Join(Directories.Data, "host.db")) {
			wg.Done()
			return
		}

		result := ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "host.db"),
			NewPath:   path.Join(Directories.Data, "host.db"),
			OldType:   oldHost{},
			NewType:   Host{},
			MigrateObject: func(o interface{}) (interface{}, error) {
				original := o.(oldHost)
				return Host{
					ID:          original.ID,
					Name:        original.Name,
					Address:     original.Address,
					Port:        original.Port,
					PSK:         original.PSK,
					Enabled:     original.Enabled,
					GroupIDs:    original.GroupIDs,
					Environment: environ.FromMap(original.Environment),
				}, nil
			},
		})
		if !result.Success {
			log.Fatal("Error migrating host table: %s", result.Error.Error())
		}
		log.Warn("Host store migration results: %+v", result)
		wg.Done()
	}()

	go func() {
		if !FileExists(path.Join(Directories.Data, "group.db")) {
			wg.Done()
			return
		}

		result := ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "group.db"),
			NewPath:   path.Join(Directories.Data, "group.db"),
			OldType:   oldGroup{},
			NewType:   Group{},
			MigrateObject: func(o interface{}) (interface{}, error) {
				original := o.(oldGroup)
				return Group{
					ID:          original.ID,
					Name:        original.Name,
					ScriptIDs:   original.ScriptIDs,
					Environment: environ.FromMap(original.Environment),
				}, nil
			},
		})
		if !result.Success {
			log.Fatal("Error migrating group table: %s", result.Error.Error())
		}
		log.Warn("Group store migration results: %+v", result)
		wg.Done()
	}()

	go func() {
		if !FileExists(path.Join(Directories.Data, "script.db")) {
			wg.Done()
			return
		}

		result := ds.Migrate(ds.MigrateParams{
			TablePath: path.Join(Directories.Data, "script.db"),
			NewPath:   path.Join(Directories.Data, "script.db"),
			OldType:   oldScript{},
			NewType:   Script{},
			MigrateObject: func(o interface{}) (interface{}, error) {
				original := o.(oldScript)
				return Script{
					ID:               original.ID,
					Name:             original.Name,
					Enabled:          original.Enabled,
					Executable:       original.Executable,
					Script:           original.Script,
					Environment:      environ.FromMap(original.Environment),
					UID:              original.UID,
					GID:              original.GID,
					WorkingDirectory: original.WorkingDirectory,
					AfterExecution:   original.AfterExecution,
				}, nil
			},
		})
		if !result.Success {
			log.Fatal("Error migrating script table: %s", result.Error.Error())
		}
		log.Warn("Script store migration results: %+v", result)
		wg.Done()
	}()

	wg.Wait()
}
