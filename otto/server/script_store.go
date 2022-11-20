package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/otto/server/environ"
)

func (s *scriptStoreObject) ScriptWithID(id string) (script *Script) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		script = s.scriptWithID(tx, id)
		return nil
	})
	return
}

func (s *scriptStoreObject) scriptWithID(tx ds.IReadTransaction, id string) *Script {
	obj, err := tx.Get(id)
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

func (s *scriptStoreObject) ScriptWithName(name string) (script *Script) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		script = s.scriptWithName(tx, name)
		return nil
	})
	return
}

func (s *scriptStoreObject) scriptWithName(tx ds.IReadTransaction, name string) *Script {
	obj, err := tx.GetUnique("Name", name)
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

func (s *scriptStoreObject) AllScripts() (scripts []Script) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		scripts = s.allScripts(tx)
		return nil
	})
	return
}

func (s *scriptStoreObject) allScripts(tx ds.IReadTransaction) []Script {
	objects, err := tx.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
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
	RunAs            RunAs
	WorkingDirectory string
	AfterExecution   string
	AttachmentIDs    []string
}

func (s *scriptStoreObject) NewScript(params newScriptParameters) (script *Script, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		script, err = s.newScript(tx, params)
		return nil
	})
	return
}

func (s *scriptStoreObject) newScript(tx ds.IReadWriteTransaction, params newScriptParameters) (*Script, *Error) {
	if s.scriptWithName(tx, params.Name) != nil {
		log.Warn("Script with name '%s' already exists", params.Name)
		return nil, ErrorUser("Script with name '%s' already exists", params.Name)
	}
	if params.AfterExecution != "" && !IsAgentAction(params.AfterExecution) {
		return nil, ErrorUser("Invalid agent action %s", params.AfterExecution)
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
		WorkingDirectory: params.WorkingDirectory,
		AfterExecution:   params.AfterExecution,
		AttachmentIDs:    params.AttachmentIDs,
	}
	if err := limits.Check(script); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Add(script); err != nil {
		log.Error("Error adding new script '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Added new script '%s'", params.Name)
	ScriptCache.Update(tx)
	return &script, nil
}

type editScriptParameters struct {
	Name             string
	Executable       string
	Script           string
	Environment      []environ.Variable
	RunAs            RunAs
	WorkingDirectory string
	AfterExecution   string
	AttachmentIDs    []string
}

func (s *scriptStoreObject) EditScript(script *Script, params editScriptParameters) (newScript *Script, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		newScript, err = s.editScript(tx, script, params)
		return nil
	})
	return
}

func (s *scriptStoreObject) editScript(tx ds.IReadWriteTransaction, script *Script, params editScriptParameters) (*Script, *Error) {
	if existingScript := s.scriptWithName(tx, params.Name); existingScript != nil && existingScript.ID != script.ID {
		log.Warn("Script with name '%s' already exists", params.Name)
		return nil, ErrorUser("Script with name '%s' already exists", params.Name)
	}
	if params.AfterExecution != "" && !IsAgentAction(params.AfterExecution) {
		return nil, ErrorUser("Invalid agent action %s", params.AfterExecution)
	}

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	script.Name = params.Name
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

	if err := tx.Update(*script); err != nil {
		log.Error("Error updating script '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updating script '%s'", params.Name)
	ScriptCache.Update(tx)
	return script, nil
}

func (s *scriptStoreObject) DeleteScript(script *Script) (err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		err = s.deleteScript(tx, script)
		return nil
	})
	return
}

func (s *scriptStoreObject) deleteScript(tx ds.IReadWriteTransaction, script *Script) *Error {
	for _, schedule := range ScheduleCache.All() {
		if schedule.ScriptID == script.ID {
			return ErrorUser("Script is used by schedule %s", schedule.Name)
		}
	}

	if err := tx.Delete(*script); err != nil {
		log.Error("Error deleting script '%s': %s", script.Name, err.Error())
		return ErrorFrom(err)
	}

	for _, id := range script.AttachmentIDs {
		if err := AttachmentStore.DeleteAttachment(id); err != nil {
			log.Error("Error deleting attachment '%s': %s", id, err.Message)
		}
	}

	GroupStore.CleanupDeadScripts(tx)
	log.Info("Deleting script '%s'", script.Name)
	ScriptCache.Update(tx)
	return nil
}
