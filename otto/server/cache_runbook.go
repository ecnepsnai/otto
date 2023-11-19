package server

import (
	"sync"

	"github.com/ecnepsnai/ds"
)

type cacheTypeRunbook struct {
	lock   *sync.RWMutex
	all    []Runbook
	byName map[string]int
	byID   map[string]int
}

// RunbookCache the runbook cache
var RunbookCache = &cacheTypeRunbook{lock: &sync.RWMutex{}}

// Update populate the runbook cache, will panic if not able to populate
func (c *cacheTypeRunbook) Update(tx ds.IReadTransaction) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Runbook{}
	c.byName = map[string]int{}
	c.byID = map[string]int{}
	runbooks := RunbookStore.allRunbooks(tx)

	c.all = runbooks
	for i, runbook := range runbooks {
		c.byName[runbook.Name] = i
		c.byID[runbook.ID] = i
	}

	log.Debug("Updated runbook cache")
	Stats.Counters.NumberRunbooks.Set(uint64(len(c.all)))
}

// All get all runbooks
func (c *cacheTypeRunbook) All() []Runbook {
	c.lock.RLock()
	defer c.lock.RUnlock()

	all := make([]Runbook, len(c.all))
	copy(all, c.all)

	return all
}

// ByName get an runbook by its name
func (c *cacheTypeRunbook) ByName(name string) *Runbook {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byName[name]
	if !k {
		return nil
	}
	return &c.all[idx]
}

// ByID get an runbook by its ID
func (c *cacheTypeRunbook) ByID(id string) *Runbook {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byID[id]
	if !k {
		return nil
	}
	return &c.all[idx]
}
