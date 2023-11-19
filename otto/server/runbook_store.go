package server

import (
	"reflect"

	"github.com/ecnepsnai/ds"
)

func (s *runbookStoreObject) AllRunbooks() (runbooks []Runbook) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		runbooks = s.allRunbooks(tx)
		return nil
	})
	return
}

func (s *runbookStoreObject) allRunbooks(tx ds.IReadTransaction) []Runbook {
	objects, err := tx.GetAll(&ds.GetOptions{Ascending: false})

	if err != nil {
		log.PError("Error getting all runbooks", map[string]interface{}{
			"error": err.Error(),
		})
		return []Runbook{}
	}
	count := len(objects)
	if count == 0 {
		return []Runbook{}
	}

	var runbooks = make([]Runbook, count)
	for i, object := range objects {
		runbook, ok := object.(Runbook)
		if !ok {
			log.PPanic("Invalid object type in runbook store", map[string]interface{}{
				"expected_type": reflect.TypeOf(Runbook{}).String(),
				"got_type":      reflect.TypeOf(object).String(),
			})
		}
		runbooks[i] = runbook
	}

	return runbooks
}

func (s *runbookStoreObject) RunbookWithID(id string) (runbook *Runbook) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		runbook = s.runbookWithID(tx, id)
		return nil
	})
	return
}

func (s *runbookStoreObject) runbookWithID(tx ds.IReadTransaction, id string) *Runbook {
	object, err := tx.Get(id)
	if err != nil {
		log.PError("Error getting runbook", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
		return nil
	}
	if object == nil {
		log.PWarn("No runbook found", map[string]interface{}{
			"id": id,
		})
		return nil
	}

	runbook, ok := object.(Runbook)
	if !ok {
		log.PPanic("Invalid object type in runbook store", map[string]interface{}{
			"expected_type": reflect.TypeOf(Runbook{}).String(),
			"got_type":      reflect.TypeOf(object).String(),
		})
	}

	return &runbook
}

func (s *runbookStoreObject) New(params Runbook) (runbook *Runbook, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		runbook, err = s.new(tx, params)
		return nil
	})
	return
}

func (s *runbookStoreObject) new(tx ds.IReadWriteTransaction, params Runbook) (*Runbook, *Error) {
	if params.ID == "" {
		params.ID = newID()
	}
	if s.runbookWithID(tx, params.ID) != nil {
		log.PWarn("Attempt to add duplicate runbook", map[string]interface{}{
			"name": params.ID,
		})
		return nil, ErrorUser("runbook with name '%s' already exists", params.ID)
	}

	highestRunLevel := ScriptRunLevelReadOnly
	for _, scriptID := range params.ScriptIDs {
		script := ScriptCache.ByID(scriptID)
		if script == nil {
			log.PError("Cannot create runbook associated with script that does not exist", map[string]interface{}{
				"script_id": scriptID,
			})
			return nil, ErrorUser("Script with ID '%s' not found", scriptID)
		}
		if script.RunLevel > highestRunLevel {
			highestRunLevel = script.RunLevel
		}
	}
	params.RunLevel = highestRunLevel

	for _, groupID := range params.GroupIDs {
		if group := GroupCache.ByID(groupID); group == nil {
			log.PError("Cannot create runbook associated with group that does not exist", map[string]interface{}{
				"group_id": groupID,
			})
			return nil, ErrorUser("Script with ID '%s' not found", groupID)
		}
	}

	if err := tx.Add(params); err != nil {
		log.PError("Error saving new runbook", map[string]interface{}{
			"id":    params.ID,
			"error": err.Error(),
		})
		return nil, ErrorFrom(err)
	}

	log.PInfo("Added new runbook", map[string]interface{}{
		"id": params.ID,
	})
	RunbookCache.Update(tx)

	return &params, nil
}

func (s *runbookStoreObject) Edit(id string, params Runbook) (runbook *Runbook, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		runbook, err = s.edit(tx, id, params)
		return nil
	})
	return
}

func (s *runbookStoreObject) edit(tx ds.IReadWriteTransaction, id string, params Runbook) (*Runbook, *Error) {
	runbook := s.runbookWithID(tx, id)
	if runbook == nil {
		log.PError("No runbook found", map[string]interface{}{
			"id": id,
		})
		return nil, ErrorUser("No runbook found")
	}

	params.ID = runbook.ID

	highestRunLevel := ScriptRunLevelReadOnly
	for _, scriptID := range params.ScriptIDs {
		script := ScriptCache.ByID(scriptID)
		if script == nil {
			log.PError("Cannot create runbook associated with script that does not exist", map[string]interface{}{
				"script_id": scriptID,
			})
			return nil, ErrorUser("Script with ID '%s' not found", scriptID)
		}
		if script.RunLevel > highestRunLevel {
			highestRunLevel = script.RunLevel
		}
	}
	params.RunLevel = highestRunLevel

	for _, groupID := range params.GroupIDs {
		if group := GroupCache.ByID(groupID); group == nil {
			log.PError("Cannot create runbook associated with group that does not exist", map[string]interface{}{
				"group_id": groupID,
			})
			return nil, ErrorUser("Script with ID '%s' not found", groupID)
		}
	}

	if err := tx.Update(params); err != nil {
		log.PError("Error updating runbook", map[string]interface{}{
			"id":    params.ID,
			"error": err.Error(),
		})
		return nil, ErrorFrom(err)
	}

	log.PInfo("Updated runbook", map[string]interface{}{
		"id": params.ID,
	})
	RunbookCache.Update(tx)

	return &params, nil
}

func (s *runbookStoreObject) Delete(id string) (runbook *Runbook, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		runbook, err = s.delete(tx, id)
		return nil
	})
	return
}

func (s *runbookStoreObject) delete(tx ds.IReadWriteTransaction, id string) (*Runbook, *Error) {
	runbook := s.runbookWithID(tx, id)
	if runbook == nil {
		return nil, ErrorUser("runbook not found")
	}

	if err := tx.Delete(*runbook); err != nil {
		log.PError("Error deleting runbook", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
		return nil, ErrorFrom(err)
	}

	log.PInfo("Deleted runbook", map[string]interface{}{
		"id": id,
	})
	RunbookCache.Update(tx)

	return runbook, nil
}
