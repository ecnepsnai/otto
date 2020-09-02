package server

import "github.com/ecnepsnai/otto/server/environ"

// Group describes a group object
type Group struct {
	ID          string `ds:"primary"`
	Name        string `ds:"unique"`
	ScriptIDs   []string
	Environment []environ.Variable
}

// HostIDs return the IDs for each host member of this group
func (g Group) HostIDs() []string {
	return GetGroupCache()[g.ID]
}

// Hosts get all hosts for this group
func (g Group) Hosts() ([]Host, *Error) {
	hosts := make([]Host, len(g.HostIDs()))
	for i, hostID := range g.HostIDs() {
		host, err := HostStore.HostWithID(hostID)
		if err != nil {
			return nil, err
		}
		hosts[i] = *host
	}
	return hosts, nil
}

// Scripts get all scripts for this group
func (g Group) Scripts() ([]Script, *Error) {
	scripts := make([]Script, len(g.ScriptIDs))
	for i, scriptID := range g.ScriptIDs {
		script, err := ScriptStore.ScriptWithID(scriptID)
		if err != nil {
			return nil, err
		}
		scripts[i] = *script
	}
	return scripts, nil
}
