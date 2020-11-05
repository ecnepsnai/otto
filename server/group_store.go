package server

import (
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/otto/server/environ"
)

func (s *groupStoreObject) GroupWithID(id string) (*Group, *Error) {
	obj, err := s.Table.Get(id)
	if err != nil {
		log.Error("Error getting group with ID '%s': %s", id, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	group, k := obj.(Group)
	if !k {
		log.Error("Object is not of type 'Group'")
		return nil, ErrorServer("incorrect type")
	}

	return &group, nil
}

func (s *groupStoreObject) GroupWithName(name string) (*Group, *Error) {
	obj, err := s.Table.GetUnique("Name", name)
	if err != nil {
		log.Error("Error getting group with name '%s': %s", name, err.Error())
		return nil, ErrorFrom(err)
	}
	if obj == nil {
		return nil, nil
	}
	group, k := obj.(Group)
	if !k {
		log.Error("Object is not of type 'Group'")
		return nil, ErrorServer("incorrect type")
	}

	return &group, nil
}

func (s *groupStoreObject) findDuplicate(name string) string {
	nameGroup, err := s.GroupWithName(name)
	if err != nil {
		return ""
	}
	if nameGroup != nil {
		return nameGroup.ID
	}

	return ""
}

func (s *groupStoreObject) AllGroups() ([]Group, *Error) {
	objs, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error getting all groups: %s", err.Error())
		return nil, ErrorFrom(err)
	}
	if objs == nil || len(objs) == 0 {
		return []Group{}, nil
	}

	groups := make([]Group, len(objs))
	for i, obj := range objs {
		group, k := obj.(Group)
		if !k {
			log.Error("Object is not of type 'Group'")
			return []Group{}, ErrorServer("incorrect type")
		}
		groups[i] = group
	}

	return groups, nil
}

type newGroupParameters struct {
	Name        string
	ScriptIDs   []string
	Environment []environ.Variable
}

func (s *groupStoreObject) NewGroup(params newGroupParameters) (*Group, *Error) {
	if s.findDuplicate(params.Name) != "" {
		log.Warn("Group with name '%s' already exists", params.Name)
		return nil, ErrorUser("Name already in use")
	}

	var enabledScripts = make([]string, len(params.ScriptIDs))
	for i, script := range params.ScriptIDs {
		s, err := ScriptStore.ScriptWithID(script)
		if err != nil {
			return nil, err
		}
		if s == nil {
			log.Warn("No script with ID '%s'", script)
			return nil, ErrorUser("No script with ID '%s'", script)
		}
		enabledScripts[i] = s.ID
	}

	group := Group{
		ID:          newID(),
		Name:        params.Name,
		ScriptIDs:   enabledScripts,
		Environment: params.Environment,
	}

	if err := s.Table.Add(group); err != nil {
		log.Error("Error adding new group '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Added new group '%s'", params.Name)
	UpdateGroupCache()
	return &group, nil
}

type editGroupParameters struct {
	Name        string
	ScriptIDs   []string
	Environment []environ.Variable
}

func (s *groupStoreObject) EditGroup(group *Group, params editGroupParameters) (*Group, *Error) {
	dupID := s.findDuplicate(params.Name)
	if dupID != "" && dupID != group.ID {
		log.Warn("Group with name '%s' already exists", params.Name)
		return nil, ErrorUser("Name already in use")
	}

	var enabledScripts = make([]string, len(params.ScriptIDs))
	for i, script := range params.ScriptIDs {
		s, err := ScriptStore.ScriptWithID(script)
		if err != nil {
			return nil, err
		}
		if s == nil {
			log.Warn("No script with ID '%s'", script)
			return nil, ErrorUser("No script with ID '%s'", script)
		}
		enabledScripts[i] = s.ID
	}

	group.Name = params.Name
	group.ScriptIDs = enabledScripts
	group.Environment = params.Environment

	if err := s.Table.Update(*group); err != nil {
		log.Error("Error updating group '%s': %s", params.Name, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updating group '%s'", params.Name)
	UpdateGroupCache()
	return group, nil
}

func (s *groupStoreObject) DeleteGroup(group *Group) *Error {
	hosts, err := group.Hosts()
	if err != nil {
		return err
	}
	if len(hosts) > 0 {
		return ErrorUser("Can't delete group with hosts")
	}

	if Options.Register.DefaultGroupID == group.ID {
		return ErrorUser("Can't delete group that is the default group for host registration")
	}

	for _, rule := range Options.Register.Rules {
		if rule.GroupID == group.ID {
			return ErrorUser("Can't delete group that is used in a host registration rule")
		}
	}

	schedules, _ := ScheduleStore.AllSchedulesForGroup(group.ID)
	if len(schedules) > 0 {
		return ErrorUser("Can't delete group that is used in a schedule")
	}

	if groups, _ := s.AllGroups(); len(groups) <= 1 {
		return ErrorUser("At least one group must exist")
	}

	if err := s.Table.Delete(*group); err != nil {
		log.Error("Error deleting group '%s': %s", group.Name, err.Error())
		return ErrorFrom(err)
	}

	heartbeatStore.CleanupHeartbeats()
	UpdateGroupCache()
	log.Info("Deleting group '%s'", group.Name)
	return nil
}

func (s *groupStoreObject) CleanupDeadScripts() *Error {
	groups, err := s.AllGroups()
	if err != nil {
		log.Error("Error getting all groups: %s", err.Message)
		return err
	}

	for _, group := range groups {
		i := len(group.ScriptIDs) - 1
		for i >= 0 {
			script, err := ScriptStore.ScriptWithID(group.ScriptIDs[i])
			if err != nil {
				return err
			}
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
