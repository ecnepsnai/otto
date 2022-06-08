package server

import "sync"

type cacheTypeScript struct {
	lock   *sync.RWMutex
	all    []Script
	byName map[string]int
	byID   map[string]int
}

// ScriptCache the script cache
var ScriptCache = &cacheTypeScript{lock: &sync.RWMutex{}}

// Update populate the script cache, will panic if not able to populate
func (c *cacheTypeScript) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Script{}
	c.byName = map[string]int{}
	c.byID = map[string]int{}
	scripts := ScriptStore.AllScripts()

	c.all = scripts
	for i, script := range scripts {
		c.byName[script.Name] = i
		c.byID[script.ID] = i
	}

	log.Debug("Updated script cache")
	Stats.Counters.NumberScripts.Set(uint64(len(c.all)))
}

// All get all scripts
func (c *cacheTypeScript) All() []Script {
	c.lock.RLock()
	defer c.lock.RUnlock()

	all := make([]Script, len(c.all))
	copy(all, c.all)

	return all
}

// ByName get an script by its name
func (c *cacheTypeScript) ByName(name string) *Script {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byName[name]
	if !k {
		return nil
	}
	return &c.all[idx]
}

// ByID get an script by its ID
func (c *cacheTypeScript) ByID(id string) *Script {
	c.lock.RLock()
	defer c.lock.RUnlock()

	idx, k := c.byID[id]
	if !k {
		return nil
	}
	return &c.all[idx]
}
