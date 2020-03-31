package server

import (
	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/scheduler"
)

var schedule *scheduler.Schedule

var scheduleDisabled = false

// ScheduleSetup start the scheduler
func ScheduleSetup() {
	schedule = scheduler.New([]scheduler.Job{
		{
			Pattern: "1 0 * * *",
			Name:    "RotateLogs",
			Exec: func() error {
				return logtic.Rotate()
			},
		},
		{
			Pattern: "/5 * * * *",
			Name:    "PingHosts",
			Exec: func() error {
				return HostStore.PingAll()
			},
		},
	})
	if !scheduleDisabled {
		go schedule.Start()
	}
}
