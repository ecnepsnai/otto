package server

import "github.com/ecnepsnai/ds"

// CacheSetup populate all caches
func CacheSetup() {
	HostStore.Table.StartRead(func(tx ds.IReadTransaction) error {
		HostCache.Update(tx)
		return nil
	})
	GroupStore.Table.StartRead(func(tx ds.IReadTransaction) error {
		GroupCache.Update(tx)
		return nil
	})
	ScriptStore.Table.StartRead(func(tx ds.IReadTransaction) error {
		ScriptCache.Update(tx)
		return nil
	})
	ScheduleStore.Table.StartRead(func(tx ds.IReadTransaction) error {
		ScheduleCache.Update(tx)
		return nil
	})
	UserStore.Table.StartRead(func(tx ds.IReadTransaction) error {
		UserCache.Update(tx)
		return nil
	})
}
