package server

import "sync"

type cacheTypeHost struct {
	lock    *sync.RWMutex
	all     []Host
	enabled []Host
	byName  map[string]int
	byID    map[string]int
}

// HostCache the host cache
var HostCache = &cacheTypeHost{lock: &sync.RWMutex{}}

// Update populate the host cache, will panic if not able to populate
func (c *cacheTypeHost) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Host{}
	c.enabled = []Host{}
	c.byName = map[string]int{}
	c.byID = map[string]int{}
	hosts := HostStore.AllHosts()

	c.all = hosts
	for i, host := range hosts {
		if host.Enabled {
			c.enabled = append(c.enabled, host)
		}
		c.byName[host.Name] = i
		c.byID[host.ID] = i
	}

	log.Debug("Updated host cache")
}

// All get all hosts
func (c *cacheTypeHost) All() []Host {
	c.lock.RLock()
	defer c.lock.RUnlock()

	all := make([]Host, len(c.all))
	for i, host := range c.all {
		all[i] = host
	}

	return all
}

// Enabled get all enabled hosts
func (c *cacheTypeHost) Enabled() []Host {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.enabled
}

// ByName get an host by its name
func (c *cacheTypeHost) ByName(name string) *Host {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byName[name]
	if !k {
		return nil
	}
	return &c.all[idx]
}

// ByID get an host by its ID
func (c *cacheTypeHost) ByID(id string) *Host {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byID[id]
	if !k {
		return nil
	}
	return &c.all[idx]
}
