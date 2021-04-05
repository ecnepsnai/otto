package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/otto/server/environ"
)

func (s *scriptStoreObject) ScriptWithID(id string) *Script {
	obj, err := s.Table.Get(id)
	if err != nil {
		log.Error("Error getting script: id='%s' error='%s'", id, err.Error())
		return nil
	}
	if obj == nil {
		return nil
	}
	script, k := obj.(Script)
	if !k {
		log.Fatal("Error getting script: id='%s' error='%s'", id, "invalid type")
	}

	return &script
}

func (s *scriptStoreObject) ScriptWithName(name string) *Script {
	obj, err := s.Table.GetUnique("Name", name)
	if err != nil {
		log.Error("Error getting script: name='%s' error='%s'", name, err.Error())
		return nil
	}
	if obj == nil {
		return nil
	}
	script, k := obj.(Script)
	if !k {
		log.Fatal("Error getting script: name='%s' error='%s'", name, "invalid type")
	}

	return &script
}

func (s *scriptStoreObject) AllScripts() []Script {
	objects, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error listing scripts: error='%s'", err.Error())
		return []Script{}
	}
	if len(objects) == 0 {
		return []Script{}
	}

	scripts := make([]Script, len(objects))
	for i, obj := range objects {
		script, k := obj.(Script)
		if !k {
			log.Fatal("Error listing scripts: error='%s'", "invalid type")
		}
		scripts[i] = script
	}

	return scripts
}

type newScriptParameters struct {
	Name             string
	Executable       string
	Script           string
	Environment      []environ.Variable
	RunAs            ScriptRunAs
	WorkingDirectory string
	AfterExecution   string
	AttachmentIDs    []string
}

func (s *scriptStoreObject) NewScript(params newScriptParameters) (*Script, *Error) {
	if s.ScriptWithName(params.Name) != nil {
		log.Warn("Script with name '%s' already exists", params.Name)
		return nil, ErrorUser("Script with name '%s' already exists", params.Name)
	}
	if params.AfterExecution != "" && !IsClientAction(params.AfterExecution) {
		return nil, ErrorUser("Invalid client action %s", params.AfterExecution)
	}

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	script := Script{
		ID:               newID(),
		Name:             params.Name,
		Executable:       params.Executable,
		Script:           params.Script,
		Environment:      params.Environment,
		RunAs:            params.RunAs,
		Enabled:          true,
		WorkingDirectory: params.WorkingDirectory,
		AfterExecution:   params.AfterExecution,
		AttachmentIDs:    params.AttachmentIDs,
	}
	if err := limits.Check(script); err != nil {
		return nil, ErrorUser(err.Error())
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
	RunAs            ScriptRunAs
	WorkingDirectory string
	AfterExecution   string
	AttachmentIDs    []string
}

func (s *scriptStoreObject) EditScript(script *Script, params editScriptParameters) (*Script, *Error) {
	if existingScript := s.ScriptWithName(params.Name); existingScript != nil && existingScript.ID != script.ID {
		log.Warn("Script with name '%s' already exists", params.Name)
		return nil, ErrorUser("Script with name '%s' already exists", params.Name)
	}
	if params.AfterExecution != "" && !IsClientAction(params.AfterExecution) {
		return nil, ErrorUser("Invalid client action %s", params.AfterExecution)
	}

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	script.Name = params.Name
	script.Enabled = params.Enabled
	script.Executable = params.Executable
	script.Script = params.Script
	script.Environment = params.Environment
	script.RunAs = params.RunAs
	script.WorkingDirectory = params.WorkingDirectory
	script.AfterExecution = params.AfterExecution
	script.AttachmentIDs = params.AttachmentIDs
	if err := limits.Check(script); err != nil {
		return nil, ErrorUser(err.Error())
	}

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

	for _, id := range script.AttachmentIDs {
		if err := AttachmentStore.DeleteAttachment(id); err != nil {
			log.Error("Error deleting attachment '%s': %s", id, err.Message)
		}
	}

	GroupStore.CleanupDeadScripts()
	log.Info("Deleting script '%s'", script.Name)
	return nil
}
