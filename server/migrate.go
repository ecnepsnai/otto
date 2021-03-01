package server

import (
	"fmt"
	"path"

	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/secutil"
)

var neededTableVersion = 9

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
		if i == 9 {
			migrate9()
		}
		i++
	}

	State.SetTableVersion(i)
}

// #9 update registration rules
func migrate9() {
	log.Debug("Start migrate 9")

	type oldRegisterRuleType struct {
		ID       string `ds:"primary"`
		Property string `ds:"index"`
		Pattern  string
		GroupID  string `ds:"index"`
	}
	type newRegisterRuleType struct {
		ID      string `ds:"primary"`
		Name    string `ds:"unique"`
		Clauses []RegisterRuleClause
		GroupID string `ds:"index"`
	}

	if !FileExists(path.Join(Directories.Data, "registerrule.db")) {
		return
	}

	results := ds.Migrate(ds.MigrateParams{
		TablePath: path.Join(Directories.Data, "registerrule.db"),
		NewPath:   path.Join(Directories.Data, "registerrule.db"),
		OldType:   oldRegisterRuleType{},
		NewType:   newRegisterRuleType{},
		MigrateObject: func(old interface{}) (interface{}, error) {
			oldRule, ok := old.(oldRegisterRuleType)
			if !ok {
				panic("Invalid type")
			}
			newRule := newRegisterRuleType{
				ID:   oldRule.ID,
				Name: fmt.Sprintf("Migrated Rule %s", secutil.RandomString(4)),
				Clauses: []RegisterRuleClause{
					{
						Property: oldRule.Property,
						Pattern:  oldRule.Pattern,
					},
				},
				GroupID: oldRule.GroupID,
			}
			return newRule, nil
		},
	})

	if results.Error != nil {
		log.Fatal("Error migrating registerrule database: %s", results.Error.Error())
	}

	log.Debug("Migrated registerrule: %+v", results)
}
