package server

import (
	"github.com/ecnepsnai/cron"
	"github.com/ecnepsnai/logtic"
)

var cronDisabled = false

// CronSetup start the cron
func CronSetup() {
	schedule, err := cron.New([]cron.Job{
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
			Name:    "CleanupAttachments",
			Exec: func() {
				AttachmentStore.Cleanup()
			},
		},
	})
	if err != nil {
		log.Fatal("Error starting up scheduled tasks: %s", err.Error())
	}
	if !cronDisabled {
		go schedule.Start()
	}
}
