package server

import (
	"github.com/ecnepsnai/cron"
	"github.com/ecnepsnai/logtic"
)

var cronDisabled = false

// CronSetup start the cron
func CronSetup() {
	schedule := cron.New([]cron.Job{
		{
			Pattern: "0 * * * *",
			Name:    "CleanupSessions",
			Exec: func() {
				SessionStore.CleanupSessions()
			},
		},
		{
			Pattern: "1 0 * * *",
			Name:    "RotateLogs",
			Exec: func() {
				logtic.Rotate()
			},
		},
		{
			Pattern: "* * * * *",
			Name:    "RunSchedules",
			Exec: func() {
				ScheduleStore.RunSchedules()
			},
		},
		{
			Pattern: "0 * * * *",
			Name:    "CleanupScriptFiles",
			Exec: func() {
				FileStore.Cleanup()
			},
		},
	})
	if !cronDisabled {
		go schedule.Start()
	}
}
