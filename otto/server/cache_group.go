package server

import "sync"

type cacheTypeGroup struct {
	lock    *sync.RWMutex
	all     []Group
	byName  map[string]int
	byID    map[string]int
	hostIDs map[string][]string
}

// GroupCache the group cache
var GroupCache = &cacheTypeGroup{lock: &sync.RWMutex{}}

// Update populate the group cache, will panic if not able to populate
func (c *cacheTypeGroup) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Group{}
	c.byName = map[string]int{}
	c.byID = map[string]int{}
	c.hostIDs = map[string][]string{}
	groups := GroupStore.AllGroups()

	c.all = groups
	for i, group := range groups {
		c.byName[group.Name] = i
		c.byID[group.ID] = i
		c.hostIDs[group.ID] = []string{}
	}
	for _, host := range HostCache.All() {
		for _, groupID := range host.GroupIDs {
			c.hostIDs[groupID] = append(c.hostIDs[groupID], host.ID)
		}
	}

	log.Debug("Updated group cache")
	Stats.Counters.NumberGroups.Set(uint64(len(c.all)))
}

// All get all groups
func (c *cacheTypeGroup) All() []Group {
	c.lock.RLock()
	defer c.lock.RUnlock()

	all := make([]Group, len(c.all))
	copy(all, c.all)

	return all
}

// ByName get an group by its name
func (c *cacheTypeGroup) ByName(name string) *Group {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byName[name]
	if !k {
		return nil
	}
	return &c.all[idx]
}

// ByID get an group by its ID
func (c *cacheTypeGroup) ByID(id string) *Group {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byID[id]
	if !k {
		return nil
	}
	return &c.all[idx]
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
