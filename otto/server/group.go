package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
)

// Group describes a group object
type Group struct {
	ID          string `ds:"primary"`
	Name        string `ds:"unique" min:"1" max:"140"`
	ScriptIDs   []string
	Environment []environ.Variable
}

// HostIDs return the IDs for each host member of this group
func (g *Group) HostIDs() []string {
	return GroupCache.HostIDs(g.ID)
}

// Hosts get all hosts for this group
func (g *Group) Hosts() ([]Host, *Error) {
	hosts := make([]Host, len(g.HostIDs()))
	for i, hostID := range g.HostIDs() {
		host := HostCache.ByID(hostID)
		hosts[i] = *host
	}
	return hosts, nil
}

// Scripts get all scripts for this group
func (g *Group) Scripts() ([]Script, *Error) {
	scripts := make([]Script, len(g.ScriptIDs))
	for i, scriptID := range g.ScriptIDs {
		script := ScriptCache.ByID(scriptID)
		if script == nil {
			log.PError("Group referrs to non-existant script", map[string]interface{}{
				"group_id":  g.ID,
				"script_id": scriptID,
			})
			group := *g
			group.ScriptIDs = append(group.ScriptIDs[:i], group.ScriptIDs[i+1:]...)
			GroupStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
				return tx.Update(group)
			})
			return group.Scripts()
		}
		scripts[i] = *script
	}
	return scripts, nil
}
