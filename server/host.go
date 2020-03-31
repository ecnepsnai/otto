package server

// Host describes a otto host
type Host struct {
	ID          string `ds:"primary"`
	Name        string `ds:"unique"`
	Address     string `ds:"unique"`
	Port        uint32
	PSK         string
	Enabled     bool `ds:"index"`
	GroupIDs    []string
	Environment map[string]string
}

// Groups return all groups for this host
func (h Host) Groups() ([]Group, *Error) {
	groups := make([]Group, len(h.GroupIDs))
	for i, groupID := range h.GroupIDs {
		group, err := GroupStore.GroupWithID(groupID)
		if err != nil {
			return nil, err
		}
		groups[i] = *group
	}
	return groups, nil
}

// ScriptEnabledGroup describes a host where a script is eanbled on it by a group
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
