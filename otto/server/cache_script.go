package server

import "sync"

type cacheTypeScript struct {
	lock    *sync.RWMutex
	all     []Script
	enabled []Script
	byName  map[string]int
	byID    map[string]int
}

// ScriptCache the script cache
var ScriptCache = &cacheTypeScript{lock: &sync.RWMutex{}}

// Update populate the script cache, will panic if not able to populate
func (c *cacheTypeScript) Update() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []Script{}
	c.enabled = []Script{}
	c.byName = map[string]int{}
	c.byID = map[string]int{}
	scripts := ScriptStore.AllScripts()

	c.all = scripts
	for i, script := range scripts {
		if script.Enabled {
			c.enabled = append(c.enabled, script)
		}
		c.byName[script.Name] = i
		c.byID[script.ID] = i
	}

	log.Debug("Updated script cache")
}

// All get all scripts
func (c *cacheTypeScript) All() []Script {
	c.lock.RLock()
	defer c.lock.RUnlock()

	all := make([]Script, len(c.all))
	for i, script := range c.all {
		all[i] = script
	}

	return all
}

// Enabled get all enabled scripts
func (c *cacheTypeScript) Enabled() []Script {
	c.lock.RLock()
	defer c.lock.RUnlock()

	enabled := make([]Script, len(c.enabled))
	for i, script := range c.enabled {
		enabled[i] = script
	}

	return enabled
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