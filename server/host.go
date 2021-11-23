package server

import (
	"time"

	"github.com/ecnepsnai/otto/server/environ"
)

// Host describes a otto host
type Host struct {
	ID          string `ds:"primary"`
	Name        string `ds:"unique" min:"1" max:"140"`
	Address     string `ds:"unique" min:"1"`
	Port        uint32
	Enabled     bool `ds:"index"`
	Trust       HostTrust
	GroupIDs    []string
	Environment []environ.Variable
}

type HostTrust struct {
	TrustedIdentity   string
	UntrustedIdentity string
	LastTrustUpdate   time.Time
}

// Groups return all groups for this host
func (h Host) Groups() ([]Group, *Error) {
	groups := make([]Group, len(h.GroupIDs))
	for i, groupID := range h.GroupIDs {
		group := GroupCache.ByID(groupID)
		groups[i] = *group
	}
	return groups, nil
}

// ScriptEnabledGroup describes a host where a script is enabled on it by a group
type ScriptEnabledGroup struct {
	ScriptID   string
	ScriptName string
	GroupID    string
	GroupName  string
}

// Scripts return all scripts for this host
func (h Host) Scripts() []ScriptEnabledGroup {
	hostScripts := []ScriptEnabledGroup{}
	groups, err := h.Groups()
	if err != nil {
		return nil
	}
	for _, group := range groups {
		scripts, err := group.Scripts()
		if err != nil {
			return nil
		}
		ehabledGroups := make([]ScriptEnabledGroup, len(scripts))
		for i, script := range scripts {
			ehabledGroups[i] = ScriptEnabledGroup{
				ScriptID:   script.ID,
				ScriptName: script.Name,
				GroupID:    group.ID,
				GroupName:  group.Name,
			}
		}
		hostScripts = append(hostScripts, ehabledGroups...)
	}
	return hostScripts
}
