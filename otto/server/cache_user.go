package server

import (
	"sync"

	"github.com/ecnepsnai/ds"
)

type cacheTypeUser struct {
	lock       *sync.RWMutex
	all        []User
	enabled    []User
	byUsername map[string]User
}

// UserCache the user cache
var UserCache = &cacheTypeUser{lock: &sync.RWMutex{}}

// Update populate the user cache, will panic if not able to populate
func (c *cacheTypeUser) Update(tx ds.IReadTransaction) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.all = []User{}
	c.enabled = []User{}
	c.byUsername = map[string]User{}
	users := UserStore.allUsers(tx)

	c.all = users
	for _, user := range users {
		if user.CanLogIn {
			c.enabled = append(c.enabled, user)
		}
		c.byUsername[user.Username] = user
	}

	log.Debug("Updated user cache")
	Stats.Counters.NumberUsers.Set(uint64(len(c.all)))
}

// All get all users
func (c *cacheTypeUser) All() []User {
	c.lock.RLock()
	defer c.lock.RUnlock()

	all := make([]User, len(c.all))
	copy(all, c.all)

	return all
}

// Enabled get all enabled users
func (c *cacheTypeUser) Enabled() []User {
	c.lock.RLock()
	defer c.lock.RUnlock()

	enabled := make([]User, len(c.enabled))
	copy(enabled, c.enabled)

	return enabled
}

// ByUsername get an user by its username
func (c *cacheTypeUser) ByUsername(username string) (User, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	a, k := c.byUsername[username]
	return a, k
}
