package server

import (
	"time"

	"github.com/ecnepsnai/cron"
	"github.com/ecnepsnai/ds"
)

func (s scheduleStoreObject) AllSchedules() ([]Schedule, *Error) {
	objs, err := s.Table.GetAll(&ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error getting all schedules: %s", err.Error())
		return nil, ErrorFrom(err)
	}
	if objs == nil || len(objs) == 0 {
		return []Schedule{}, nil
	}

	schedules := make([]Schedule, len(objs))
	for i, obj := range objs {
		host, k := obj.(Schedule)
		if !k {
			log.Error("Object is not of type 'Schedule'")
			return []Schedule{}, ErrorServer("incorrect type")
		}
		schedules[i] = host
	}

	return schedules, nil
}

func (s scheduleStoreObject) AllSchedulesForScript(scriptID string) ([]Schedule, *Error) {
	objs, err := s.Table.GetIndex("ScriptID", scriptID, &ds.GetOptions{Sorted: true, Ascending: true})
	if err != nil {
		log.Error("Error getting all schedules for script %s: %s", scriptID, err.Error())
		return nil, ErrorFrom(err)
	}
	if objs == nil || len(objs) == 0 {
		return []Schedule{}, nil
	}

	schedules := make([]Schedule, len(objs))
	for i, obj := range objs {
		host, k := obj.(Schedule)
		if !k {
			log.Error("Object is not of type 'Schedule'")
			return []Schedule{}, ErrorServer("incorrect type")
		}
		schedules[i] = host
	}

	return schedules, nil
}

func (s scheduleStoreObject) AllSchedulesForGroup(groupID string) ([]Schedule, *Error) {
	matchedSchedules := []Schedule{}
	schedules, err := s.AllSchedules()
	if err != nil {
		return nil, err
	}
	for _, schedule := range schedules {
		if StringSliceContains(groupID, schedule.Scope.GroupIDs) {
			matchedSchedules = append(matchedSchedules, schedule)
		}
	}

	return schedules, nil
}

func (s scheduleStoreObject) AllSchedulesForHost(hostID string) ([]Schedule, *Error) {
	matchedSchedules := []Schedule{}
	schedules, err := s.AllSchedules()
	if err != nil {
		return nil, err
	}
	for _, schedule := range schedules {
		if len(schedule.Scope.GroupIDs) > 0 {
			for _, groupID := range schedule.Scope.GroupIDs {
				hostIDs, ok := GetGroupCache()[groupID]
				if !ok {
					continue
				}
				for _, h := range hostIDs {
					if hostID == h {
						matchedSchedules = append(matchedSchedules, schedule)
						break
					}
				}
			}
		} else if len(schedule.Scope.HostIDs) > 0 {
			if StringSliceContains(hostID, schedule.Scope.HostIDs) {
				matchedSchedules = append(matchedSchedules, schedule)
			}
		}
	}

	return schedules, nil
}

func (s scheduleStoreObject) RunSchedules() {
	schedules, err := s.AllSchedules()
	if err != nil {
		log.Error("Error fetching all schedules: %s", err.Message)
		return
	}

	for _, schedule := range schedules {
		if !schedule.Enabled {
			log.Debug("Skipping disabled schedule: %s", schedule.ID)
			continue
		}

		j := cron.Job{Pattern: schedule.Pattern}
		if j.WouldRunNow() {
			schedule.RunNow()
		}
	}
}

type newScheduleParameters struct {
	ScriptID string `ds:"index"`
	Scope    ScheduleScope
	Pattern  string
}

func (s *scheduleStoreObject) NewSchedule(params newScheduleParameters) (*Schedule, *Error) {
	if len(params.Scope.GroupIDs) > 0 && len(params.Scope.HostIDs) > 0 {
		return nil, ErrorUser("Cannot specify both group IDs and host IDs")
	}
	if len(params.Scope.GroupIDs) <= 0 && len(params.Scope.HostIDs) <= 0 {
		return nil, ErrorUser("Must specify at least one group or host")
	}

	if script, _ := ScriptStore.ScriptWithID(params.ScriptID); script == nil {
		return nil, ErrorUser("Unknown script ID '%s'", params.ScriptID)
	}

	for _, groupID := range params.Scope.GroupIDs {
		if group, _ := GroupStore.GroupWithID(groupID); group == nil {
			return nil, ErrorUser("Unknown group ID '%s'", groupID)
		}
	}

	for _, hostID := range params.Scope.HostIDs {
		if host, _ := HostStore.HostWithID(hostID); host == nil {
			return nil, ErrorUser("Unknown host ID '%s'", hostID)
		}
	}

	schedule := Schedule{
		ID:       NewID(),
		ScriptID: params.ScriptID,
		Scope: ScheduleScope{
			HostIDs:  params.Scope.HostIDs,
			GroupIDs: params.Scope.GroupIDs,
		},
		Pattern: params.Pattern,
		Enabled: true,
	}

	if err := s.Table.Add(schedule); err != nil {
		log.Error("Error adding new schedule '%s': %s", schedule.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Added new schedule '%s'", schedule.ID)
	return &schedule, nil
}

type editScheduleParameters struct {
	Scope   ScheduleScope
	Pattern string
	Enabled bool
}

func (s *scheduleStoreObject) EditSchedule(schedule *Schedule, params editScheduleParameters) (*Schedule, *Error) {
	if len(params.Scope.GroupIDs) > 0 && len(params.Scope.HostIDs) > 0 {
		return nil, ErrorUser("Cannot specify both group IDs and host IDs")
	}
	if len(params.Scope.GroupIDs) <= 0 && len(params.Scope.HostIDs) <= 0 {
		return nil, ErrorUser("Must specify at least one group or host")
	}

	for _, groupID := range params.Scope.GroupIDs {
		if group, _ := GroupStore.GroupWithID(groupID); group == nil {
			return nil, ErrorUser("Unknown group ID '%s'", groupID)
		}
	}

	for _, hostID := range params.Scope.HostIDs {
		if host, _ := HostStore.HostWithID(hostID); host == nil {
			return nil, ErrorUser("Unknown host ID '%s'", hostID)
		}
	}

	schedule.Scope.HostIDs = params.Scope.HostIDs
	schedule.Scope.GroupIDs = params.Scope.GroupIDs
	schedule.Pattern = params.Pattern
	schedule.Enabled = params.Enabled

	if err := s.Table.Update(*schedule); err != nil {
		log.Error("Error updating schedule '%s': %s", schedule.ID, err.Error())
		return nil, ErrorFrom(err)
	}

	log.Info("Updated schedule '%s'", schedule.ID)
	return schedule, nil
}

func (s *scheduleStoreObject) DeleteSchedule(schedule *Schedule) *Error {
	if err := s.Table.Delete(*schedule); err != nil {
		log.Error("Error deleting schedule '%s': %s", schedule.ID, err.Error())
		return ErrorFrom(err)
	}

	log.Info("Deleted schedule '%s'", schedule.ID)
	return nil
}

func (s *scheduleStoreObject) updateLastRun(schedule Schedule) *Error {
	schedule.LastRunTime = time.Now()
	if err := s.Table.Update(schedule); err != nil {
		log.Error("Error updating last run for schedule '%s': %s", schedule.ID, err.Error())
		return ErrorFrom(err)
	}
	return nil
}
