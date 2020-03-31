package server

import (
	"sync"
)

// WarmCache warm all caches
func WarmCache() {
	UpdateGroupCache()
}

var groupCacheLock = &sync.RWMutex{}
var groupCacheCurrent = map[string][]string{}

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
	groups, err := GroupStore.AllGroups()
	if err != nil {
		log.Fatal("Error building group cache: %s", err.Message)
	}
	for _, group := range groups {
		groupCacheCurrent[group.ID] = []string{}
	}

	hosts, err := HostStore.AllHosts()
	if err != nil {
		log.Fatal("Error building group cache: %s", err.Message)
	}
	for _, host := range hosts {
		for _, groupID := range host.GroupIDs {
			groupCacheCurrent[groupID] = append(groupCacheCurrent[groupID], host.ID)
		}
	}
}
