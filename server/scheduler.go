package server

import (
	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/scheduler"
)

var schedulerDisabled = false

// SchedulerSetup start the scheduler
func SchedulerSetup() {
	schedule := scheduler.New([]scheduler.Job{
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
	if !schedulerDisabled {
		go schedule.Start()
	}
}
