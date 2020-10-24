package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
)

// Script describes an otto script
type Script struct {
	ID               string `ds:"primary"`
	Name             string `ds:"unique"`
	Enabled          bool   `ds:"index"`
	Executable       string
	Script           string
	Environment      []environ.Variable
	UID              uint32
	GID              uint32
	WorkingDirectory string
	AfterExecution   string
	FileIDs          []string
}

func (s *scriptStoreObject) ScriptWithID(id string) (*Script, *Error) {
	obj, err := s.Table.Get(id)
	if err != nil {
		log.Error("Error getting script with ID '%s': %s", id, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	script, k := obj.(Script)
	if !k {
		log.Error("Object is not of type 'Script'")
		return nil, ErrorServer("incorrect type")
	}

	return &script, nil
}

func (s *scriptStoreObject) ScriptWithName(name string) (*Script, *Error) {
	obj, err := s.Table.GetUnique("Name", name)
	if err != nil {
		log.Error("Error getting script with name '%s': %s", name, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	script, k := obj.(Script)
	if !k {
		log.Error("Object is not of type 'Script'")
		return nil, ErrorServer("incorrect type")
	}

	return &script, nil
}

func (s *scriptStoreObject) AllScripts() ([]Script, *Error) {
	objs, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error getting all scripts: %s", err.Error())
		return nil, ErrorFrom(err)
	}
	if objs == nil || len(objs) == 0 {
		return []Script{}, nil
	}

	scripts := make([]Script, len(objs))
	for i, obj := range objs {
		script, k := obj.(Script)
		if !k {
			log.Error("Object is not of type 'Script'")
			return []Script{}, ErrorServer("incorrect type")
		}
		scripts[i] = script
	}

	return scripts, nil
}

type newScriptParameters struct {
	Name             string
	Executable       string
	Script           string
	Environment      []environ.Variable
	UID              uint32
	GID              uint32
	WorkingDirectory string
	AfterExecution   string
	FileIDs          []string
}

func (s *scriptStoreObject) NewScript(params newScriptParameters) (*Script, *Error) {
	existingScript, err := s.ScriptWithName(params.Name)
	if err != nil {
		return nil, err
	}
	if existingScript != nil {
		log.Warn("Script with name '%s' already exists", params.Name)
		return nil, ErrorUser("Script with name '%s' already exists", params.Name)
	}
	if params.AfterExecution != "" && !IsClientAction(params.AfterExecution) {
		return nil, ErrorUser("Invalid client action %s", params.AfterExecution)
	}

	script := Script{
		ID:               NewID(),
		Name:             params.Name,
		Executable:       params.Executable,
		Script:           params.Script,
		Environment:      params.Environment,
		UID:              params.UID,
		GID:              params.GID,
		Enabled:          true,
		WorkingDirectory: params.WorkingDirectory,
		AfterExecution:   params.AfterExecution,
		FileIDs:          params.FileIDs,
	}

	if err := s.Table.Add(script); err != nil {
		log.Error("Error adding new script '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Added new script '%s'", params.Name)
	return &script, nil
}

type editScriptParameters struct {
	Name             string
	Enabled          bool
	Executable       string
	Script           string
	Environment      []environ.Variable
	UID              uint32
	GID              uint32
	WorkingDirectory string
	AfterExecution   string
	FileIDs          []string
}

func (s *scriptStoreObject) EditScript(script *Script, params editScriptParameters) (*Script, *Error) {
	existingScript, err := s.ScriptWithName(params.Name)
	if err != nil {
		return nil, err
	}
	if existingScript != nil && existingScript.ID != script.ID {
		log.Warn("Script with name '%s' already exists", params.Name)
		return nil, ErrorUser("Script with name '%s' already exists", params.Name)
	}
	if params.AfterExecution != "" && !IsClientAction(params.AfterExecution) {
		return nil, ErrorUser("Invalid client action %s", params.AfterExecution)
	}

	script.Name = params.Name
	script.Enabled = params.Enabled
	script.Executable = params.Executable
	script.Script = params.Script
	script.Environment = params.Environment
	script.UID = params.UID
	script.GID = params.GID
	script.WorkingDirectory = params.WorkingDirectory
	script.AfterExecution = params.AfterExecution
	script.FileIDs = params.FileIDs

	if err := s.Table.Update(*script); err != nil {
		log.Error("Error updating script '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updating script '%s'", params.Name)
	return script, nil
}

func (s *scriptStoreObject) DeleteScript(script *Script) *Error {
	if err := s.Table.Delete(*script); err != nil {
		log.Error("Error deleting script '%s': %s", script.Name, err.Error())
		return ErrorFrom(err)
	}

	for _, id := range script.FileIDs {
		if err := FileStore.DeleteFile(id); err != nil {
			log.Error("Error deleting script file '%s': %s", id, err.Message)
		}
	}

	GroupStore.CleanupDeadScripts()
	log.Info("Deleting script '%s'", script.Name)
	return nil
}

// Groups all groups with this script enabled
func (s *Script) Groups() []Group {
	enabledGroups := []Group{}

	groups, err := GroupStore.AllGroups()
	if err != nil {
		return []Group{}
	}
	for _, group := range groups {
		hasScript := false
		for _, scriptID := range group.ScriptIDs {
			if scriptID == s.ID {
				hasScript = true
				break
			}
		}
		if !hasScript {
			continue
		}
		enabledGroups = append(enabledGroups, group)
	}

	return enabledGroups
}

// ScriptEnabledHost describes a host where a script is eanbled on it by a group
type ScriptEnabledHost struct {
	ScriptID   string
	ScriptName string
	GroupID    string
	GroupName  string
	HostID     string
	HostName   string
}

// Hosts all hosts with this script enabled
func (s *Script) Hosts() []ScriptEnabledHost {
	enabledHosts := []ScriptEnabledHost{}

	for _, group := range s.Groups() {
		hosts, err := group.Hosts()
		if err != nil {
			return []ScriptEnabledHost{}
		}
		ehs := make([]ScriptEnabledHost, len(hosts))
		for i, host := range hosts {
			ehs[i] = ScriptEnabledHost{
				ScriptID:   s.ID,
				ScriptName: s.Name,
				GroupID:    group.ID,
				GroupName:  group.Name,
				HostID:     host.ID,
				HostName:   host.Name,
			}
		}
		enabledHosts = append(enabledHosts, ehs...)
	}

	return enabledHosts
}

func (s *scriptStoreObject) SetGroups(script *Script, groupIDs []string) *Error {
	groups := map[string]bool{}
	allGroups, err := GroupStore.AllGroups()
	if err != nil {
		return err
	}
	for _, group := range allGroups {
		var i = -1
		for y, groupID := range groupIDs {
			if groupID == group.ID {
				i = y
				break
			}
		}
		groups[group.ID] = i != -1
	}

	for groupID, enable := range groups {
		group, err := GroupStore.GroupWithID(groupID)
		if err != nil {
			return err
		}
		if group == nil {
			return ErrorUser("No group with ID %s", groupID)
		}

		var i = -1
		for y, scriptID := range group.ScriptIDs {
			if scriptID == script.ID {
				i = y
				break
			}
		}

		if i == -1 && enable {
			group.ScriptIDs = append(group.ScriptIDs, script.ID)
			log.Debug("Enabling script '%s' on group '%s'", script.Name, group.Name)
		} else if i != -1 && !enable {
			group.ScriptIDs = append(group.ScriptIDs[:i], group.ScriptIDs[i+1:]...)
			log.Debug("Disabling script '%s' on group '%s'", script.Name, group.Name)
		} else {
			continue
		}

		if err := GroupStore.Table.Update(*group); err != nil {
			log.Error("Error updating group '%s': %s", group.Name, err.Error())
			return ErrorFrom(err)
		}
	}

	return nil
}

// Files all files for this script
func (s *Script) Files() ([]File, *Error) {
	if len(s.FileIDs) == 0 {
		return []File{}, nil
	}

	files := make([]File, len(s.FileIDs))
	for i, id := range s.FileIDs {
		file, err := FileStore.FileWithID(id)
		if err != nil {
			return nil, err
		}
		if file == nil {
			log.Error("File '%s' does not exist, found on script '%s'", id, s.ID)
			return nil, ErrorServer("missing file")
		}
		files[i] = *file
	}

	return files, nil
}
