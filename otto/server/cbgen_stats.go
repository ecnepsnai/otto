package server

// This file is was generated automatically by Codegen v1.12.0
// Do not make changes to this file as they will be lost

import (
	"time"

	"github.com/ecnepsnai/stats"
)

type cbgenStatsCounters struct {
	NumberGroups     *stats.Counter
	NumberHosts      *stats.Counter
	NumberSchedules  *stats.Counter
	NumberScripts    *stats.Counter
	NumberUsers      *stats.Counter
	ReachableHosts   *stats.Counter
	TrustedHosts     *stats.Counter
	UnreachableHosts *stats.Counter
	UntrustedHosts   *stats.Counter
}

type cbgenStatsTimedCounters struct {
}

type cbgenStatsTimers struct {
}

type cbgenStatsObject struct {
	Counters      cbgenStatsCounters
	TimedCounters cbgenStatsTimedCounters
	Timers        cbgenStatsTimers
}

// Stats the global stats object
var Stats *cbgenStatsObject

// statsSetup setup the stats object
func statsSetup() {
	Stats = &cbgenStatsObject{
		Counters: cbgenStatsCounters{
			NumberGroups:     stats.NewCounter(),
			NumberHosts:      stats.NewCounter(),
			NumberSchedules:  stats.NewCounter(),
			NumberScripts:    stats.NewCounter(),
			NumberUsers:      stats.NewCounter(),
			ReachableHosts:   stats.NewCounter(),
			TrustedHosts:     stats.NewCounter(),
			UnreachableHosts: stats.NewCounter(),
			UntrustedHosts:   stats.NewCounter(),
		},
		TimedCounters: cbgenStatsTimedCounters{},
		Timers:        cbgenStatsTimers{},
	}
}

// Reset reset all volatile stats
func (s *cbgenStatsObject) Reset() {
	statsSetup()
}

// GetCounterValues get a map of current counters
func (s *cbgenStatsObject) GetCounterValues() map[string]uint64 {
	return map[string]uint64{
		"NumberGroups":     s.Counters.NumberGroups.Get(),
		"NumberHosts":      s.Counters.NumberHosts.Get(),
		"NumberSchedules":  s.Counters.NumberSchedules.Get(),
		"NumberScripts":    s.Counters.NumberScripts.Get(),
		"NumberUsers":      s.Counters.NumberUsers.Get(),
		"ReachableHosts":   s.Counters.ReachableHosts.Get(),
		"TrustedHosts":     s.Counters.TrustedHosts.Get(),
		"UnreachableHosts": s.Counters.UnreachableHosts.Get(),
		"UntrustedHosts":   s.Counters.UntrustedHosts.Get(),
	}
}

// GetTimedCounterValues get a map of all timed counter values
func (s *cbgenStatsObject) GetTimedCounterValues() map[string]uint64 {
	return map[string]uint64{}
}

// GetTimedCounterValuesFrom get a map of all timed counter values
func (s *cbgenStatsObject) GetTimedCounterValuesFrom(d time.Duration) map[string]uint64 {
	return map[string]uint64{}
}

// GetTimerAverages get the average times for all timers
func (s *cbgenStatsObject) GetTimerAverages() map[string]time.Duration {
	return map[string]time.Duration{}
}

// GetTimerValues get all vaues for all timers
func (s *cbgenStatsObject) GetTimerValues() map[string][]time.Duration {
	return map[string][]time.Duration{}
}
