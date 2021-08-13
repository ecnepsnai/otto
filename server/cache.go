package server

// CacheSetup populate all caches
func CacheSetup() {
	HostCache.Update()
	GroupCache.Update()
	ScriptCache.Update()
	ScheduleCache.Update()
	UserCache.Update()
}
