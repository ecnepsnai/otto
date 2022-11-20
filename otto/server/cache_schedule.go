package server

import (
	"sync"

	"github.com/ecnepsnai/ds"
)

type cacheTypeSchedule struct {
	lock    *sync.RWMutex
	all     []Schedule
	enabled []Schedule
	byName  map[string]int
	byID    map[string]int
}

// ScheduleCache the schedule cache
var ScheduleCache = &cacheTypeSchedule{lock: &sync.RWMutex{}}

// Update populate the schedule cache, will panic if not able to populate
func (c *cacheTypeSchedule) Update(tx ds.IReadTransaction) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Schedule{}
	c.enabled = []Schedule{}
	c.byName = map[string]int{}
	c.byID = map[string]int{}
	schedules := ScheduleStore.allSchedules(tx)

	c.all = schedules
	for i, schedule := range schedules {
		if schedule.Enabled {
			c.enabled = append(c.enabled, schedule)
		}
		c.byName[schedule.Name] = i
		c.byID[schedule.ID] = i
	}

	log.Debug("Updated schedule cache")
	Stats.Counters.NumberSchedules.Set(uint64(len(c.all)))
}

// All get all schedules
func (c *cacheTypeSchedule) All() []Schedule {
	c.lock.RLock()
	defer c.lock.RUnlock()

	all := make([]Schedule, len(c.all))
	copy(all, c.all)

	return all
}

// Enabled get all enabled schedules
func (c *cacheTypeSchedule) Enabled() []Schedule {
	c.lock.RLock()
	defer c.lock.RUnlock()

	enabled := make([]Schedule, len(c.enabled))
	copy(enabled, c.enabled)

	return enabled
}

// ByName get an schedule by its name
func (c *cacheTypeSchedule) ByName(name string) *Schedule {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byName[name]
	if !k {
		return nil
	}
	return &c.all[idx]
}

// ByID get an schedule by its ID
func (c *cacheTypeSchedule) ByID(id string) *Schedule {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byID[id]
	if !k {
		return nil
	}
	schedule := c.all[idx]
	return &schedule
}
