// Package freq is a set of structures for forming a task schedule within an accuracy of 1 minute
package freq

import (
	"fmt"
	"time"
)

// EveryUnit enum type for the unit to match the Every property on a Frequency object
type EveryUnit int

const (
	// EveryHour enum value for every hour
	EveryHour = EveryUnit(1)
	// EveryDay enum value for every day
	EveryDay = EveryUnit(2)
	// EveryWeek enum value for every week
	EveryWeek = EveryUnit(3)
	// EveryMonth enum value for every month
	EveryMonth = EveryUnit(4)
)

// Frequency describes the frequency that a task should run
type Frequency struct {
	// Every how frequent the task should run. Use a predefined enum value.
	Every EveryUnit
	// AtTime the time at which the task should run. Date values (year, month, day) are ignored.
	AtTime time.Time
}

func addHours(t time.Time, h int) time.Time {
	return time.Unix(t.Unix()+int64(3600*h), 0)
}

func addMinutes(t time.Time, m int) time.Time {
	return time.Unix(t.Unix()+int64(60*m), 0)
}

func (f Frequency) nextRunEveryHour(lastRunTime time.Time) time.Time {
	minutesSinceLastRun := time.Since(lastRunTime).Minutes()
	minutes := 60 - int(minutesSinceLastRun)
	return addMinutes(time.Now(), minutes)
}

func (f Frequency) nextRunEveryDay(lastRunTime time.Time) time.Time {
	hoursSinceLastRun := time.Since(lastRunTime).Hours()
	hours := 24 - int(hoursSinceLastRun)

	n := time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		f.AtTime.Hour(),
		f.AtTime.Minute(),
		f.AtTime.Second(),
		0,
		time.Now().Location(),
	)

	return addHours(n, hours)
}

func (f Frequency) nextRunEveryWeek(lastRunTime time.Time) time.Time {
	daysSinceLastRun := time.Since(lastRunTime).Hours() / 24.0
	days := 7 - int(daysSinceLastRun)

	n := time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		f.AtTime.Hour(),
		f.AtTime.Minute(),
		f.AtTime.Second(),
		0,
		time.Now().Location(),
	)

	return n.AddDate(0, 0, days)
}

func (f Frequency) nextRunEveryMonth(lastRunTime time.Time) time.Time {
	firstOfCurrentMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	firstOfNextMonth := time.Now().AddDate(0, 1, -time.Now().Day()+1)
	daysInMonth := int(firstOfNextMonth.Sub(firstOfCurrentMonth).Hours() / 24.0)
	fmt.Printf("Days in month: %d\n", daysInMonth)

	daysSinceLastRun := time.Since(lastRunTime).Hours() / 24.0
	fmt.Printf("Days since last run: %f\n", daysSinceLastRun)
	days := daysInMonth - int(daysSinceLastRun)

	n := time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		f.AtTime.Hour(),
		f.AtTime.Minute(),
		f.AtTime.Second(),
		0,
		time.Now().Location(),
	)

	fmt.Printf("Adding %d days\n", days)
	return n.AddDate(0, 0, days)
}

// NextRunTime get the next time this task should run as determined by the last time it ran. If the task has never ran,
// provide time.Unix(0, 0)
func (f Frequency) NextRunTime(lastRunTime time.Time) time.Time {
	switch f.Every {
	case EveryHour:
		return f.nextRunEveryHour(lastRunTime)
	case EveryDay:
		return f.nextRunEveryDay(lastRunTime)
	case EveryWeek:
		return f.nextRunEveryWeek(lastRunTime)
	case EveryMonth:
		return f.nextRunEveryMonth(lastRunTime)
	}

	panic("unknown Every value")
}
