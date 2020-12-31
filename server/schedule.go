package server

import (
	"time"

	"github.com/ecnepsnai/set"
)

// Schedule describes a recurring task
type Schedule struct {
	ID          string `ds:"primary"`
	Name        string `ds:"unique"`
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
	for i, id := range s.GroupIDs {
		g, err := GroupStore.GroupWithID(id)
		if err != nil {
			return nil, err
		}
		groups[i] = *g
	}

	return groups, nil
}

// Hosts get the hosts for this schedule
func (s ScheduleScope) Hosts() ([]Host, *Error) {
	hostIDs := set.NewString()
	if len(s.GroupIDs) > 0 {
		for _, id := range s.GroupIDs {
			groupHosts, ok := GetGroupCache()[id]
			if !ok {
				return nil, ErrorServer("empty cache")
			}
			for _, hostID := range groupHosts {
				hostIDs.Add(hostID)
			}
		}
	} else if len(s.HostIDs) > 0 {
		for _, hostID := range s.HostIDs {
			hostIDs.Add(hostID)
		}
	}

	if hostIDs.Length() == 0 {
		log.Warn("Schedule with no hosts or groups: %s")
		return []Host{}, nil
	}

	hosts := make([]Host, hostIDs.Length())
	for i, id := range hostIDs.Values() {
		g, err := HostStore.HostWithID(id)
		if err != nil {
			return nil, err
		}
		hosts[i] = *g
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
			hostIDs, ok := GetGroupCache()[id]
			if !ok {
				log.Error("Group cache empty, cannot run scheduled script")
				return
			}
			for _, hostID := range hostIDs {
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
		host, err := HostStore.HostWithID(hostID)

		script, err := ScriptStore.ScriptWithID(s.ScriptID)
		if err != nil {
			continue
		}
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
