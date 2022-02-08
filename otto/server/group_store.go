package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
	"github.com/ecnepsnai/otto/server/environ"
)

func (s *groupStoreObject) GroupWithID(id string) *Group {
	object, err := s.Table.Get(id)
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

func (s *groupStoreObject) GroupWithName(name string) *Group {
	object, err := s.Table.GetUnique("Name", name)
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

func (s *groupStoreObject) AllGroups() []Group {
	objects, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
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

func (s *groupStoreObject) NewGroup(params newGroupParameters) (*Group, *Error) {
	if s.GroupWithName(params.Name) != nil {
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

	if err := s.Table.Add(group); err != nil {
		log.Error("Error adding new group '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Added new group '%s'", params.Name)
	GroupCache.Update()
	return &group, nil
}

type editGroupParameters struct {
	Name        string
	ScriptIDs   []string
	Environment []environ.Variable
}

func (s *groupStoreObject) EditGroup(group *Group, params editGroupParameters) (*Group, *Error) {
	if existingGroup := s.GroupWithName(params.Name); existingGroup != nil && existingGroup.ID != group.ID {
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

	if err := s.Table.Update(*group); err != nil {
		log.Error("Error updating group '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updating group '%s'", params.Name)
	GroupCache.Update()
	return group, nil
}

func (s *groupStoreObject) DeleteGroup(group *Group) *Error {
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

	if groups := s.AllGroups(); len(groups) <= 1 {
		log.Error("At least one group must exist")
		return ErrorUser("At least one group must exist")
	}

	if err := s.Table.Delete(*group); err != nil {
		log.Error("Error deleting group '%s': %s", group.Name, err.Error())
		return ErrorFrom(err)
	}

	heartbeatStore.CleanupHeartbeats()
	GroupCache.Update()
	log.Info("Deleting group '%s'", group.Name)
	return nil
}

func (s *groupStoreObject) CleanupDeadScripts() *Error {
	groups := s.AllGroups()
	for _, group := range groups {
		i := len(group.ScriptIDs) - 1
		for i >= 0 {
			script := ScriptStore.ScriptWithID(group.ScriptIDs[i])
			if script == nil {
				log.Warn("Removing non-existant script '%s' from group '%s'", group.ScriptIDs[i], group.Name)
				group.ScriptIDs = append(group.ScriptIDs[:i], group.ScriptIDs[i+1:]...)
			}
			i--
		}
		s.Table.Update(group)
	}

	return nil
}
