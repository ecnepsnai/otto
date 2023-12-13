package server

// This file is was generated automatically by GenGo v1.13.0
// Do not make changes to this file as they will be lost

import (
	"time"

	"github.com/ecnepsnai/stats"
)

type gengoStatsCounters struct {
	NumberGroups     *stats.Counter
	NumberHosts      *stats.Counter
	NumberRunbooks   *stats.Counter
	NumberSchedules  *stats.Counter
	NumberScripts    *stats.Counter
	NumberUsers      *stats.Counter
	ReachableHosts   *stats.Counter
	TrustedHosts     *stats.Counter
	UnreachableHosts *stats.Counter
	UntrustedHosts   *stats.Counter
}

type gengoStatsTimedCounters struct {
}

type gengoStatsTimers struct {
}

type gengoStatsObject struct {
	Counters      gengoStatsCounters
	TimedCounters gengoStatsTimedCounters
	Timers        gengoStatsTimers
}

// Stats the global stats object
var Stats *gengoStatsObject

// statsSetup setup the stats object
func statsSetup() {
	Stats = &gengoStatsObject{
		Counters: gengoStatsCounters{
			NumberGroups:     stats.NewCounter(),
			NumberHosts:      stats.NewCounter(),
			NumberRunbooks:   stats.NewCounter(),
			NumberSchedules:  stats.NewCounter(),
			NumberScripts:    stats.NewCounter(),
			NumberUsers:      stats.NewCounter(),
			ReachableHosts:   stats.NewCounter(),
			TrustedHosts:     stats.NewCounter(),
			UnreachableHosts: stats.NewCounter(),
			UntrustedHosts:   stats.NewCounter(),
		},
		TimedCounters: gengoStatsTimedCounters{},
		Timers:        gengoStatsTimers{},
	}
}

// Reset reset all volatile stats
func (s *gengoStatsObject) Reset() {
	statsSetup()
}

// GetCounterValues get a map of current counters
func (s *gengoStatsObject) GetCounterValues() map[string]uint64 {
	return map[string]uint64{
		"NumberGroups":     s.Counters.NumberGroups.Get(),
		"NumberHosts":      s.Counters.NumberHosts.Get(),
		"NumberRunbooks":   s.Counters.NumberRunbooks.Get(),
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
func (s *gengoStatsObject) GetTimedCounterValues() map[string]uint64 {
	return map[string]uint64{}
}

// GetTimedCounterValuesFrom get a map of all timed counter values
func (s *gengoStatsObject) GetTimedCounterValuesFrom(d time.Duration) map[string]uint64 {
	return map[string]uint64{}
}

// GetTimerAverages get the average times for all timers
func (s *gengoStatsObject) GetTimerAverages() map[string]time.Duration {
	return map[string]time.Duration{}
}

// GetTimerValues get all vaues for all timers
func (s *gengoStatsObject) GetTimerValues() map[string][]time.Duration {
	return map[string][]time.Duration{}
}
