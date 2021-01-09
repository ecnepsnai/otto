package server

import (
	"encoding/json"
	"os"
	"path"
)

var neededTableVersion = 7

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
		if i == 7 {
			migrate7()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #7 Migrate registration rules
func migrate7() {
	log.Debug("Start migrate 7")

	type oldRegisterRule struct {
		Property string
		Pattern  string
		GroupID  string
	}

	type oldOttoOptionsRegister struct {
		Rules []oldRegisterRule
	}

	type oldOttoOptions struct {
		Register oldOttoOptionsRegister
	}

	configPath := path.Join(Directories.Data, configFileName)
	if !FileExists(configPath) {
		return
	}

	oldOptions := oldOttoOptions{}
	f, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Error opening config file '%s': %s", configPath, err.Error())
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&oldOptions); err != nil {
		log.Fatal("Error opening config file '%s': %s", configPath, err.Error())
	}

	if len(oldOptions.Register.Rules) <= 0 {
		return
	}

	cbgenDataStoreRegisterGroupStore()
	cbgenDataStoreRegisterRegisterRuleStore()
	for _, oldRule := range oldOptions.Register.Rules {
		newRule := newRegisterRuleParams{
			Property: oldRule.Property,
			Pattern:  oldRule.Pattern,
			GroupID:  oldRule.GroupID,
		}
		if newRule.Property == "uname" {
			newRule.Property = RegisterRulePropertyKernelName
		}

		if _, err := RegisterRuleStore.NewRule(newRule); err != nil {
			log.Error("Error migrating register rule to new format. Old rule: %+v. Error: %s", oldRule, err)
		}
	}
	RegisterRuleStore.Table.Close()
	GroupStore.Table.Close()
}
