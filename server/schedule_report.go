package server

import (
	"time"

	"github.com/ecnepsnai/ds"
)

// ScheduleReport describes a report from a schedule
type ScheduleReport struct {
	ID         string `ds:"primary"`
	ScheduleID string `ds:"index"`
	HostIDs    []string
	Time       ScheduleReportTime
	Result     int
}

// ScheduleReportTime describes timing information from a schedule run
type ScheduleReportTime struct {
	Start          time.Time
	Finished       time.Time
	ElapsedSeconds float64
}

func (s *schedulereportStoreObject) GetReportsForSchedule(scheduleID string) []ScheduleReport {
	objs, err := s.Table.GetIndex("ScheduleID", scheduleID, &ds.GetOptions{Sorted: true})
	if err != nil {
		log.Error("Error getting all schedule reports: %s", err.Error())
		return []ScheduleReport{}
	}
	if objs == nil || len(objs) == 0 {
		return []ScheduleReport{}
	}

	reports := make([]ScheduleReport, len(objs))
	for i, obj := range objs {
		host, k := obj.(ScheduleReport)
		if !k {
			log.Error("Object is not of type 'ScheduleReport'")
			return []ScheduleReport{}
		}
		reports[i] = host
	}

	return reports
}
