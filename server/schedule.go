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
			Pattern: "0 * * * *",
			Name:    "CleanupSessions",
			Exec: func() error {
				return SessionStore.CleanupSessions().Error
			},
		},
		{
			Pattern: "1 0 * * *",
			Name:    "RotateLogs",
			Exec: func() error {
				return logtic.Rotate()
			},
		},
	})
	if !scheduleDisabled {
		go schedule.Start()
	}
}
