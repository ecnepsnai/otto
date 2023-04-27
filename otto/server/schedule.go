package server

import (
	"time"

	"github.com/ecnepsnai/ds"
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
	groups := []Group{}
	for _, groupID := range s.GroupIDs {
		group := GroupCache.ByID(groupID)
		if group == nil {
			log.Warn("Schedule contains unknown group %s", groupID)
			continue
		}
		groups = append(groups, *group)
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

	hosts := []Host{}
	for _, hostID := range hostIDs.Values() {
		host := HostCache.ByID(hostID)
		if host == nil {
			log.Warn("Schedule contains unknown host %s", hostID)
			continue
		}
		hosts = append(hosts, *host)
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

	script := ScriptCache.ByID(s.ScriptID)
	if script == nil {
		log.PError("Schedule targets non-existant script", map[string]interface{}{
			"schedule_id": s.ID,
			"script_id":   s.ScriptID,
		})
		return
	}

	report.HostIDs = hosts.Values()
	report.HostResult = map[string]int{}
	success := 0
	fail := 0

	for _, hostID := range hosts.Values() {
		host := HostCache.ByID(hostID)
		if host == nil {
			log.PError("Schedule targets non-existant host", map[string]interface{}{
				"schedule_id": s.ID,
				"host_id":     hostID,
			})
			continue
		}

		result, err := host.RunScript(script, nil, nil)
		if err != nil {
			fail++
			log.PError("Error running scheduled script", map[string]interface{}{
				"schedule_id": s.ID,
				"script_id":   s.ScriptID,
				"host_id":     host.ID,
				"error":       err.Error(),
			})
			report.HostResult[host.ID] = 1
			continue
		} else {
			success++
			report.HostResult[host.ID] = result.Result.Code
		}

		EventStore.ScriptRun(script, host, &result.Result, &s, "")
		log.PInfo("Finished running scheduled script", map[string]interface{}{
			"schedule_id": s.ID,
			"script_id":   s.ScriptID,
			"host_id":     host.ID,
		})
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

	ScheduleReportStore.Table.StartWrite(func(tx ds.IReadWriteTransaction) error {
		return tx.Add(report)
	})
	log.PInfo("Finished running schedule", map[string]interface{}{
		"schedule_id":   s.ID,
		"start_time":    start,
		"finished_time": finished,
		"elapsed":       time.Since(start).String(),
		"num_success":   success,
		"num_fail":      fail,
	})
}
