package server

import (
	"sync"
)

// CacheSetup populate all caches
func CacheSetup() {
	HostCache.Update()
	GroupCache.Update()
	ScriptCache.Update()
	ScheduleCache.Update()
	UserCache.Update()
}

type cacheTypeHost struct {
	lock    *sync.RWMutex
	all     []Host
	enabled []Host
	byName  map[string]Host
}

// HostCache the host cache
var HostCache = &cacheTypeHost{lock: &sync.RWMutex{}}

// Update populate the host cache, will panic if not able to populate
func (c *cacheTypeHost) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Host{}
	c.enabled = []Host{}
	c.byName = map[string]Host{}
	hosts := HostStore.AllHosts()

	c.all = hosts
	for _, host := range hosts {
		if host.Enabled {
			c.enabled = append(c.enabled, host)
		}
		c.byName[host.Name] = host
	}

	log.Info("Updated host cache")
}

// All get all hosts
func (c *cacheTypeHost) All() []Host {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.all
}

// Enabled get all enabled hosts
func (c *cacheTypeHost) Enabled() []Host {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.enabled
}

// ByName get an host by its name
func (c *cacheTypeHost) ByName(name string) (Host, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	a, k := c.byName[name]
	return a, k
}

type cacheTypeGroup struct {
	lock    *sync.RWMutex
	all     []Group
	byName  map[string]Group
	hostIDs map[string][]string
}

// GroupCache the group cache
var GroupCache = &cacheTypeGroup{lock: &sync.RWMutex{}}

// Update populate the group cache, will panic if not able to populate
func (c *cacheTypeGroup) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Group{}
	c.byName = map[string]Group{}
	c.hostIDs = map[string][]string{}
	groups := GroupStore.AllGroups()

	c.all = groups
	for _, group := range groups {
		c.byName[group.Name] = group
		c.hostIDs[group.ID] = []string{}
	}
	for _, host := range HostCache.All() {
		for _, groupID := range host.GroupIDs {
			c.hostIDs[groupID] = append(c.hostIDs[groupID], host.ID)
		}
	}

	log.Info("Updated group cache")
}

// All get all groups
func (c *cacheTypeGroup) All() []Group {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.all
}

// ByName get an group by its name
func (c *cacheTypeGroup) ByName(name string) (Group, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	a, k := c.byName[name]
	return a, k
}

// Membership return a mapping of group IDs to a list of host IDs belong to that group
func (c *cacheTypeGroup) Membership() map[string][]string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.hostIDs
}

// HostIDs return host IDs for the given group
func (c *cacheTypeGroup) HostIDs(id string) []string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	ids, ok := c.hostIDs[id]
	if !ok {
		return []string{}
	}

	return ids
}

type cacheTypeScript struct {
	lock    *sync.RWMutex
	all     []Script
	enabled []Script
	byName  map[string]Script
}

// ScriptCache the script cache
var ScriptCache = &cacheTypeScript{lock: &sync.RWMutex{}}

// Update populate the script cache, will panic if not able to populate
func (c *cacheTypeScript) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Script{}
	c.enabled = []Script{}
	c.byName = map[string]Script{}
	scripts := ScriptStore.AllScripts()

	c.all = scripts
	for _, script := range scripts {
		if script.Enabled {
			c.enabled = append(c.enabled, script)
		}
		c.byName[script.Name] = script
	}

	log.Info("Updated script cache")
}

// All get all scripts
func (c *cacheTypeScript) All() []Script {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.all
}

// Enabled get all enabled scripts
func (c *cacheTypeScript) Enabled() []Script {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.enabled
}

// ByName get an script by its name
func (c *cacheTypeScript) ByName(name string) (Script, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	a, k := c.byName[name]
	return a, k
}

type cacheTypeSchedule struct {
	lock    *sync.RWMutex
	all     []Schedule
	enabled []Schedule
	byName  map[string]Schedule
}

// ScheduleCache the schedule cache
var ScheduleCache = &cacheTypeSchedule{lock: &sync.RWMutex{}}

// Update populate the schedule cache, will panic if not able to populate
func (c *cacheTypeSchedule) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Schedule{}
	c.enabled = []Schedule{}
	c.byName = map[string]Schedule{}
	schedules := ScheduleStore.AllSchedules()

	c.all = schedules
	for _, schedule := range schedules {
		if schedule.Enabled {
			c.enabled = append(c.enabled, schedule)
		}
		c.byName[schedule.Name] = schedule
	}

	log.Info("Updated schedule cache")
}

// All get all schedules
func (c *cacheTypeSchedule) All() []Schedule {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.all
}

// Enabled get all enabled schedules
func (c *cacheTypeSchedule) Enabled() []Schedule {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.enabled
}

// ByName get an schedule by its name
func (c *cacheTypeSchedule) ByName(name string) (Schedule, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	a, k := c.byName[name]
	return a, k
}

type cacheTypeUser struct {
	lock       *sync.RWMutex
	all        []User
	enabled    []User
	byUsername map[string]User
	byEmail    map[string]User
}

// UserCache the user cache
var UserCache = &cacheTypeUser{lock: &sync.RWMutex{}}

// Update populate the user cache, will panic if not able to populate
func (c *cacheTypeUser) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []User{}
	c.enabled = []User{}
	c.byUsername = map[string]User{}
	c.byEmail = map[string]User{}
	users := UserStore.AllUsers()

	c.all = users
	for _, user := range users {
		if user.CanLogIn {
			c.enabled = append(c.enabled, user)
		}
		c.byUsername[user.Username] = user
		c.byEmail[user.Email] = user
	}

	log.Info("Updated user cache")
}

// All get all users
func (c *cacheTypeUser) All() []User {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.all
}

// Enabled get all enabled users
func (c *cacheTypeUser) Enabled() []User {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.enabled
}

// ByUsername get an user by its username
func (c *cacheTypeUser) ByUsername(username string) (User, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	a, k := c.byUsername[username]
	return a, k
}

// ByEmail get an user by its email
func (c *cacheTypeUser) ByEmail(username string) (User, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	a, k := c.byEmail[username]
	return a, k
}
