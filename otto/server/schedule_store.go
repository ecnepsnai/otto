package server

import (
	"time"

	"github.com/ecnepsnai/cron"
	"github.com/ecnepsnai/ds"
	"github.com/ecnepsnai/limits"
)

func (s *scheduleStoreObject) AllSchedules() (schedules []Schedule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		schedules = s.allSchedules(tx)
		return nil
	})
	return
}

func (s *scheduleStoreObject) allSchedules(tx ds.IReadTransaction) []Schedule {
	objects, err := tx.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error listing all schedules: error='%s'", err.Error())
		return []Schedule{}
	}
	if len(objects) == 0 {
		return []Schedule{}
	}

	schedules := make([]Schedule, len(objects))
	for i, obj := range objects {
		host, k := obj.(Schedule)
		if !k {
			log.Fatal("Error listing all schedules: error='%s'", "invalid type")
		}
		schedules[i] = host
	}

	return schedules
}

func (s *scheduleStoreObject) AllSchedulesForScript(scriptID string) (schedules []Schedule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		schedules = s.allSchedulesForScript(tx, scriptID)
		return nil
	})
	return
}

func (s *scheduleStoreObject) allSchedulesForScript(tx ds.IReadTransaction, scriptID string) []Schedule {
	objects, err := tx.GetIndex("ScriptID", scriptID, &ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error listing all schedules for script: script_id='%s' error='%s'", scriptID, err.Error())
		return []Schedule{}
	}
	if len(objects) == 0 {
		return []Schedule{}
	}

	schedules := make([]Schedule, len(objects))
	for i, obj := range objects {
		host, k := obj.(Schedule)
		if !k {
			log.Fatal("Error listing all schedules for script: script_id='%s' error='%s'", scriptID, "invalid type")
		}
		schedules[i] = host
	}

	return schedules
}

func (s *scheduleStoreObject) AllSchedulesForGroup(groupID string) (schedules []Schedule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		schedules = s.allSchedulesForGroup(tx, groupID)
		return nil
	})
	return
}

func (s *scheduleStoreObject) allSchedulesForGroup(tx ds.IReadTransaction, groupID string) []Schedule {
	matchedSchedules := []Schedule{}
	schedules := s.allSchedules(tx)
	for _, schedule := range schedules {
		if stringSliceContains(groupID, schedule.Scope.GroupIDs) {
			matchedSchedules = append(matchedSchedules, schedule)
		}
	}

	return matchedSchedules
}

func (s *scheduleStoreObject) AllSchedulesForHost(hostID string) (schedules []Schedule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		schedules = s.allSchedulesForHost(tx, hostID)
		return nil
	})
	return
}

func (s *scheduleStoreObject) allSchedulesForHost(tx ds.IReadTransaction, hostID string) []Schedule {
	matchedSchedules := []Schedule{}
	schedules := s.allSchedules(tx)
	for _, schedule := range schedules {
		if len(schedule.Scope.GroupIDs) > 0 {
			for _, groupID := range schedule.Scope.GroupIDs {
				for _, h := range GroupCache.HostIDs(groupID) {
					if hostID == h {
						matchedSchedules = append(matchedSchedules, schedule)
						break
					}
				}
			}
		} else if len(schedule.Scope.HostIDs) > 0 {
			if stringSliceContains(hostID, schedule.Scope.HostIDs) {
				matchedSchedules = append(matchedSchedules, schedule)
			}
		}
	}

	return matchedSchedules
}

func (s *scheduleStoreObject) ScheduleWithID(id string) (schedule *Schedule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		schedule = s.scheduleWithID(tx, id)
		return nil
	})
	return
}

func (s *scheduleStoreObject) scheduleWithID(tx ds.IReadTransaction, id string) *Schedule {
	object, err := tx.Get(id)
	if err != nil {
		log.Error("Error getting schedule: id='%s' error='%s'", id, err.Error())
		return nil
	}
	if object == nil {
		return nil
	}

	schedule, ok := object.(Schedule)
	if !ok {
		log.Fatal("Error getting schedule: id='%s' error='%s'", id, "invalid type")
	}
	return &schedule
}

func (s *scheduleStoreObject) ScheduleWithName(name string) (schedule *Schedule) {
	s.Table.StartRead(func(tx ds.IReadTransaction) error {
		schedule = s.scheduleWithName(tx, name)
		return nil
	})
	return
}

func (s *scheduleStoreObject) scheduleWithName(tx ds.IReadTransaction, name string) *Schedule {
	object, err := tx.GetUnique("Name", name)
	if err != nil {
		log.Error("Error getting schedule: name='%s' error='%s'", name, err.Error())
		return nil
	}
	if object == nil {
		return nil
	}

	schedule, ok := object.(Schedule)
	if !ok {
		log.Fatal("Error getting schedule: name='%s' error='%s'", name, "invalid type")
	}
	return &schedule
}

func (s *scheduleStoreObject) RunSchedules() {
	schedules := s.AllSchedules()
	for _, schedule := range schedules {
		if !schedule.Enabled {
			log.Debug("Skipping disabled schedule: %s", schedule.ID)
			continue
		}

		j := cron.Job{Pattern: schedule.Pattern}
		if j.WouldRunNowInTZ(time.UTC) {
			schedule.RunNow()
		}
	}
}

type newScheduleParameters struct {
	ScriptID string
	Name     string
	Scope    ScheduleScope
	Pattern  string
}

func (s *scheduleStoreObject) NewSchedule(params newScheduleParameters) (schedule *Schedule, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		schedule, err = s.newSchedule(tx, params)
		return nil
	})
	return
}

func (s *scheduleStoreObject) newSchedule(tx ds.IReadWriteTransaction, params newScheduleParameters) (*Schedule, *Error) {
	if len(params.Scope.GroupIDs) > 0 && len(params.Scope.HostIDs) > 0 {
		return nil, ErrorUser("Cannot specify both group IDs and host IDs")
	}
	if len(params.Scope.GroupIDs) <= 0 && len(params.Scope.HostIDs) <= 0 {
		return nil, ErrorUser("Must specify at least one group or host")
	}
	if schedule, _ := tx.GetUnique("Name", params.Name); schedule != nil {
		return nil, ErrorUser("Duplicate script name")
	}
	if script := ScriptStore.ScriptWithID(params.ScriptID); script == nil {
		return nil, ErrorUser("Unknown script ID '%s'", params.ScriptID)
	}

	for _, groupID := range params.Scope.GroupIDs {
		if group := GroupCache.ByID(groupID); group == nil {
			return nil, ErrorUser("Unknown group ID '%s'", groupID)
		}
	}

	for _, hostID := range params.Scope.HostIDs {
		if host := HostCache.ByID(hostID); host == nil {
			return nil, ErrorUser("Unknown host ID '%s'", hostID)
		}
	}

	schedule := Schedule{
		ID:       newID(),
		Name:     params.Name,
		ScriptID: params.ScriptID,
		Scope: ScheduleScope{
			HostIDs:  params.Scope.HostIDs,
			GroupIDs: params.Scope.GroupIDs,
		},
		Pattern: params.Pattern,
		Enabled: true,
	}
	if err := limits.Check(schedule); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Add(schedule); err != nil {
		log.Error("Error adding new schedule '%s': %s", schedule.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Added new schedule '%s'", schedule.ID)
	ScheduleCache.Update(tx)

	return &schedule, nil
}

type editScheduleParameters struct {
	Name    string
	Scope   ScheduleScope
	Pattern string
	Enabled bool
}

func (s *scheduleStoreObject) EditSchedule(schedule *Schedule, params editScheduleParameters) (newSchedule *Schedule, err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		newSchedule, err = s.editSchedule(tx, schedule, params)
		return nil
	})
	return
}

func (s *scheduleStoreObject) editSchedule(tx ds.IReadWriteTransaction, schedule *Schedule, params editScheduleParameters) (*Schedule, *Error) {
	if len(params.Scope.GroupIDs) > 0 && len(params.Scope.HostIDs) > 0 {
		return nil, ErrorUser("Cannot specify both group IDs and host IDs")
	}
	if len(params.Scope.GroupIDs) <= 0 && len(params.Scope.HostIDs) <= 0 {
		return nil, ErrorUser("Must specify at least one group or host")
	}
	if existing := s.scheduleWithName(tx, params.Name); existing != nil && existing.ID != schedule.ID {
		log.PWarn("Schedule rename collission", map[string]interface{}{
			"schedule_id":   schedule.ID,
			"existing_id":   existing.ID,
			"schedule_name": schedule.Name,
			"existing_name": existing.Name,
		})
		return nil, ErrorUser("Schedule with name '%s' already exists", params.Name)
	}

	for _, groupID := range params.Scope.GroupIDs {
		if group := GroupCache.ByID(groupID); group == nil {
			return nil, ErrorUser("Unknown group ID '%s'", groupID)
		}
	}

	for _, hostID := range params.Scope.HostIDs {
		if host := HostCache.ByID(hostID); host == nil {
			return nil, ErrorUser("Unknown host ID '%s'", hostID)
		}
	}

	schedule.Name = params.Name
	schedule.Scope.HostIDs = params.Scope.HostIDs
	schedule.Scope.GroupIDs = params.Scope.GroupIDs
	schedule.Pattern = params.Pattern
	schedule.Enabled = params.Enabled
	if err := limits.Check(schedule); err != nil {
		return nil, ErrorUser(err.Error())
	}

	if err := tx.Update(*schedule); err != nil {
		log.Error("Error updating schedule '%s': %s", schedule.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updated schedule '%s'", schedule.ID)
	ScheduleCache.Update(tx)

	return schedule, nil
}

func (s *scheduleStoreObject) DeleteSchedule(schedule *Schedule) (err *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		err = s.deleteSchedule(tx, schedule)
		return nil
	})
	return
}

func (s *scheduleStoreObject) deleteSchedule(tx ds.IReadWriteTransaction, schedule *Schedule) *Error {
	if err := tx.Delete(*schedule); err != nil {
		log.Error("Error deleting schedule '%s': %s", schedule.ID, err.Error())
		return ErrorFrom(err)
	}

	log.Info("Deleted schedule '%s'", schedule.ID)
	ScheduleCache.Update(tx)

	return nil
}

func (s *scheduleStoreObject) updateLastRun(schedule Schedule) (rerr *Error) {
	s.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		schedule.LastRunTime = time.Now()
		if err := tx.Update(schedule); err != nil {
			log.Error("Error updating last run for schedule '%s': %s", schedule.ID, err.Error())
			rerr = ErrorFrom(err)
			return nil
		}
		ScheduleCache.Update(tx)
		return nil
	})
	return
}
