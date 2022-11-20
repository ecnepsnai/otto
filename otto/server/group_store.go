package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/otto/server/environ"
)

func (s *groupStoreObject) GroupWithID(id string) (group *Group) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		group = s.groupWithID(tx, id)
		return nil
	})
	return
}

func (s *groupStoreObject) groupWithID(tx ds.IReadTransaction, id string) *Group {
	object, err := tx.Get(id)
	if err != nil {
		log.Error("Error getting group: id='%s' error='%s'", id, err.Error())
		return nil
	}
	if object == nil {
		return nil
	}
	group, k := object.(Group)
	if !k {
		log.Fatal("Error getting group: id='%s' error='%s'", id, "invalid type")
	}

	return &group
}

func (s *groupStoreObject) GroupWithName(name string) (group *Group) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		group = s.groupWithName(tx, name)
		return nil
	})
	return
}

func (s *groupStoreObject) groupWithName(tx ds.IReadTransaction, name string) *Group {
	object, err := tx.GetUnique("Name", name)
	if err != nil {
		log.Error("Error getting group: name='%s' error='%s'", name, err.Error())
		return nil
	}
	if object == nil {
		return nil
	}
	group, k := object.(Group)
	if !k {
		log.Fatal("Error getting group: name='%s' error='%s'", name, "invalid type")
	}

	return &group
}

func (s *groupStoreObject) AllGroups() (groups []Group) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		groups = s.allGroups(tx)
		return nil
	})
	return
}

func (s *groupStoreObject) allGroups(tx ds.IReadTransaction) []Group {
	objects, err := tx.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error listing all groups: error='%s'", err.Error())
		return []Group{}
	}
	if len(objects) == 0 {
		return []Group{}
	}

	groups := make([]Group, len(objects))
	for i, obj := range objects {
		group, k := obj.(Group)
		if !k {
			log.Fatal("Error listing all groups: error='%s'", "invalid type")
		}
		groups[i] = group
	}

	return groups
}

type newGroupParameters struct {
	Name        string
	ScriptIDs   []string
	Environment []environ.Variable
}

func (s *groupStoreObject) NewGroup(params newGroupParameters) (group *Group, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		group, err = s.newGroup(tx, params)
		return nil
	})
	return
}

func (s *groupStoreObject) newGroup(tx ds.IReadWriteTransaction, params newGroupParameters) (*Group, *Error) {
	if s.groupWithName(tx, params.Name) != nil {
		log.Warn("Group with name '%s' already exists", params.Name)
		return nil, ErrorUser("Name already in use")
	}

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	var enabledScripts = make([]string, len(params.ScriptIDs))
	for i, scriptID := range params.ScriptIDs {
		script := ScriptStore.ScriptWithID(scriptID)
		if script == nil {
			log.Warn("No script with ID '%s'", scriptID)
			return nil, ErrorUser("No script with ID '%s'", scriptID)
		}
		enabledScripts[i] = script.ID
	}

	group := Group{
		ID:          newID(),
		Name:        params.Name,
		ScriptIDs:   enabledScripts,
		Environment: params.Environment,
	}
	if err := limits.Check(group); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Add(group); err != nil {
		log.Error("Error adding new group '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Added new group '%s'", params.Name)
	GroupCache.Update(tx)
	return &group, nil
}

type editGroupParameters struct {
	Name        string
	ScriptIDs   []string
	Environment []environ.Variable
}

func (s *groupStoreObject) EditGroup(group *Group, params editGroupParameters) (newGroup *Group, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		newGroup, err = s.editGroup(tx, group, params)
		return nil
	})
	return
}

func (s *groupStoreObject) editGroup(tx ds.IReadWriteTransaction, group *Group, params editGroupParameters) (*Group, *Error) {
	if existingGroup := s.groupWithName(tx, params.Name); existingGroup != nil && existingGroup.ID != group.ID {
		log.Warn("Group with name '%s' already exists", params.Name)
		return nil, ErrorUser("Name already in use")
	}

	if err := environ.Validate(params.Environment); err != nil {
		return nil, ErrorUser(err.Error())
	}

	var enabledScripts = make([]string, len(params.ScriptIDs))
	for i, scriptID := range params.ScriptIDs {
		script := ScriptStore.ScriptWithID(scriptID)
		if script == nil {
			log.Warn("No script with ID '%s'", scriptID)
			return nil, ErrorUser("No script with ID '%s'", scriptID)
		}
		enabledScripts[i] = script.ID
	}

	group.Name = params.Name
	group.ScriptIDs = enabledScripts
	group.Environment = params.Environment
	if err := limits.Check(group); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Update(*group); err != nil {
		log.Error("Error updating group '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updating group '%s'", params.Name)
	GroupCache.Update(tx)
	return group, nil
}

func (s *groupStoreObject) DeleteGroup(group *Group) (err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		err = s.deleteGroup(tx, group)
		return nil
	})
	return
}

func (s *groupStoreObject) deleteGroup(tx ds.IReadWriteTransaction, group *Group) *Error {
	hosts, err := group.Hosts()
	if err != nil {
		return err
	}
	if len(hosts) > 0 {
		log.Error("Can't delete group '%s' with hosts", group.Name)
		return ErrorUser("Can't delete group with hosts")
	}

	if Options.Register.DefaultGroupID == group.ID {
		log.Error("Can't delete group '%s' that is the default group for host registration", group.Name)
		return ErrorUser("Can't delete group that is the default group for host registration")
	}

	if rules := RegisterRuleStore.RulesForGroup(group.ID); len(rules) > 0 {
		log.Error("Can't delete group '%s' that is used in a host registration rule", group.Name)
		return ErrorUser("Can't delete group that is used in a host registration rule")
	}

	schedules := ScheduleStore.AllSchedulesForGroup(group.ID)
	if len(schedules) > 0 {
		log.Error("Can't delete group '%s' that is used in a schedule", group.Name)
		return ErrorUser("Can't delete group that is used in a schedule")
	}

	if groups := s.allGroups(tx); len(groups) <= 1 {
		log.Error("At least one group must exist")
		return ErrorUser("At least one group must exist")
	}

	if err := tx.Delete(*group); err != nil {
		log.Error("Error deleting group '%s': %s", group.Name, err.Error())
		return ErrorFrom(err)
	}

	HostStore.Table.StartRead(func(hostTx ds.IReadTransaction) error {
		heartbeatStore.CleanupHeartbeats(hostTx)
		return nil
	})
	GroupCache.Update(tx)
	log.Info("Deleting group '%s'", group.Name)
	return nil
}

func (s *groupStoreObject) CleanupDeadScripts(scriptStoreTx ds.IReadTransaction) *Error {
	err := s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		groups := s.allGroups(tx)
		for _, group := range groups {
			i := len(group.ScriptIDs) - 1
			for i >= 0 {
				script := ScriptStore.scriptWithID(scriptStoreTx, group.ScriptIDs[i])
				if script == nil {
					log.Warn("Removing non-existant script '%s' from group '%s'", group.ScriptIDs[i], group.Name)
					group.ScriptIDs = append(group.ScriptIDs[:i], group.ScriptIDs[i+1:]...)
				}
				i--
			}
			if err := tx.Update(group); err != nil {
				return err
			}
		}
		return nil
	})
	return ErrorFrom(err)
}
