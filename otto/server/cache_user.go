package server

import "sync"

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

// ByEmail get an user by its email
func (c *cacheTypeUser) ByEmail(username string) (User, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	a, k := c.byEmail[username]
	return a, k
}
