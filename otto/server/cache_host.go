package server

import (
	"sync"

	"github.com/ecnepsnai/ds"
)

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
func (c *cacheTypeHost) Update(tx ds.IReadTransaction) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Host{}
	c.enabled = []Host{}
	c.byName = map[string]int{}
	c.byID = map[string]int{}
	hosts := HostStore.allHosts(tx)

	nTrusted := uint64(0)
	nUntrusted := uint64(0)

	c.all = hosts
	for i, host := range hosts {
		if host.Enabled {
			c.enabled = append(c.enabled, host)
		}
		c.byName[host.Name] = i
		c.byID[host.ID] = i
		if host.Trust.TrustedIdentity != "" {
			nTrusted++
		} else {
			nUntrusted++
		}
	}

	log.Debug("Updated host cache")
	Stats.Counters.NumberHosts.Set(uint64(len(c.all)))
	Stats.Counters.TrustedHosts.Set(nTrusted)
	Stats.Counters.UntrustedHosts.Set(nUntrusted)
}

// All get all hosts
func (c *cacheTypeHost) All() []Host {
	c.lock.RLock()
	defer c.lock.RUnlock()

	all := make([]Host, len(c.all))
	copy(all, c.all)

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
