package server

import (
	"time"

	"github.com/ecnepsnai/set"
)

// Schedule describes a recurring task
type Schedule struct {
	ID          string `ds:"primary"`
	Name        string `ds:"unique" min:"1" max:"140"`
	ScriptID    string `ds:"index"`
	Scope       ScheduleScope
	Pattern     string
	Enabled     bool
	LastRunTime time.Time
}

// ScheduleScope describes the scope for a schedule
type ScheduleScope struct {
	HostIDs  []string
	GroupIDs []string
}

// Groups get the groups for this schedule
func (s ScheduleScope) Groups() ([]Group, *Error) {
	groups := make([]Group, len(s.GroupIDs))
	for i, groupID := range s.GroupIDs {
		group := GroupStore.GroupWithID(groupID)
		groups[i] = *group
	}

	return groups, nil
}

// Hosts get the hosts for this schedule
func (s ScheduleScope) Hosts() ([]Host, *Error) {
	hostIDs := set.NewString()
	if len(s.GroupIDs) > 0 {
		for _, id := range s.GroupIDs {
			for _, hostID := range GroupCache.HostIDs(id) {
				hostIDs.Add(hostID)
			}
		}
	} else if len(s.HostIDs) > 0 {
		for _, hostID := range s.HostIDs {
			hostIDs.Add(hostID)
		}
	}

	if hostIDs.Length() == 0 {
		log.Warn("Schedule with no hosts or groups")
		return []Host{}, nil
	}

	hosts := make([]Host, hostIDs.Length())
	for i, hostID := range hostIDs.Values() {
		host := HostStore.HostWithID(hostID)
		if host == nil {
			log.Warn("Schedule contains unknown host '%s'", hostID)
			continue
		}
		hosts[i] = *host
	}

	return hosts, nil
}

// RunNow run the schedule now
func (s Schedule) RunNow() {
	if !s.Enabled {
		return
	}

	report := ScheduleReport{
		ID:         newID(),
		ScheduleID: s.ID,
	}
	start := time.Now()

	hosts := set.NewString()
	if len(s.Scope.GroupIDs) > 0 {
		for _, id := range s.Scope.GroupIDs {
			for _, hostID := range GroupCache.HostIDs(id) {
				hosts.Add(hostID)
			}
		}
	} else if len(s.Scope.HostIDs) > 0 {
		for _, hostID := range s.Scope.HostIDs {
			hosts.Add(hostID)
		}
	}

	report.HostIDs = hosts.Values()
	report.HostResult = map[string]int{}
	success := 0
	fail := 0

	for _, hostID := range hosts.Values() {
		host := HostStore.HostWithID(hostID)
		if host == nil {
			log.Error("Schedule for nonexistant host: schedule=%s host=%s", s.ID, hostID)
			return
		}
		script := ScriptStore.ScriptWithID(s.ScriptID)
		if script == nil {
			log.Error("Schedule for nonexistant script: schedule=%s script=%s", s.ID, s.ScriptID)
			return
		}

		result, err := host.RunScript(script, nil, nil)
		if err != nil {
			fail++
			log.Error("Error running scheduled script: schedule=%s script=%s host=%s error='%s'", s.ID, s.ScriptID, host.ID, err.Message)
			report.HostResult[host.ID] = -1
			continue
		} else {
			report.HostResult[host.ID] = result.Result.Code
			success++
		}

		EventStore.ScriptRun(script, host, &result.Result, &s, "")
		log.Info("Result: %v", result)
	}

	ScheduleStore.updateLastRun(s)

	finished := time.Now()
	report.Time = ScheduleReportTime{
		Start:          start,
		Finished:       finished,
		ElapsedSeconds: time.Since(start).Seconds(),
	}
	if fail == 0 {
		report.Result = ScheduleResultSuccess
	} else if success > 0 {
		report.Result = ScheduleResultPartialSuccess
	} else {
		report.Result = ScheduleResultFail
	}
	go ScheduleReportStore.Table.Add(report)
}
