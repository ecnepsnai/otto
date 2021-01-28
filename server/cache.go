package server

import (
	"sync"
)

// WarmCache warm all caches
func WarmCache() {
	UpdateGroupCache()
	UpdateUserCache()
}

var groupCacheLock = &sync.RWMutex{}
var groupCacheCurrent = map[string][]string{}

var userCacheLock = &sync.RWMutex{}
var userCacheCurrent = map[string]User{}

// GetGroupCache get the group cache
func GetGroupCache() map[string][]string {
	groupCacheLock.RLock()
	defer groupCacheLock.RUnlock()
	return groupCacheCurrent
}

// UpdateGroupCache update the group cache
func UpdateGroupCache() {
	groupCacheLock.Lock()
	defer groupCacheLock.Unlock()

	groupCacheCurrent = map[string][]string{}
	groups := GroupStore.AllGroups()
	for _, group := range groups {
		groupCacheCurrent[group.ID] = []string{}
	}

	hosts := HostStore.AllHosts()
	for _, host := range hosts {
		for _, groupID := range host.GroupIDs {
			groupCacheCurrent[groupID] = append(groupCacheCurrent[groupID], host.ID)
		}
	}
}

// GetUserCache get the user cache
func GetUserCache() map[string]User {
	userCacheLock.RLock()
	defer userCacheLock.RUnlock()
	return userCacheCurrent
}

// UpdateUserCache update the user cache
func UpdateUserCache() {
	userCacheLock.Lock()
	defer userCacheLock.Unlock()

	userCacheCurrent = map[string]User{}
	users := UserStore.AllUsers()
	for _, user := range users {
		userCacheCurrent[user.Username] = user
	}
}
