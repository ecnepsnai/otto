package freq_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ecnepsnai/otto/server/freq"
)

var midnightLastWeek = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -7)
var midnightYesterday = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -1)
var midnightToday = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
var midnightTomorrow = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 1)
var midnightNextWeek = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, 7)
var halfHourAgo = time.Unix(time.Now().Unix()-1800, 0)
var oneHourAgo = time.Unix(time.Now().Unix()-3600, 0)
var halfHourFromNow = time.Unix(time.Now().Unix()+1800, 0)
var oneHourFromNow = time.Unix(time.Now().Unix()+3600, 0)
var oneDayFromNow = time.Now().AddDate(0, 0, 1)
var oneWeekFromNow = time.Now().AddDate(0, 0, 7)
var oneMonthFromNow = time.Now().AddDate(0, 1, 0)

func isInPast(l time.Time) bool {
	return time.Since(l).Seconds() > 0
}

func dateMatches(l, r time.Time) bool {
	// Checks if two dates match within the accuracy of 1 second
	difference := l.Unix() - r.Unix()
	return difference == 0
}

func TestEveryHour(t *testing.T) {
	f := freq.Frequency{
		Every:  freq.EveryHour,
		AtTime: midnightLastWeek,
	}

	next := f.NextRunTime(oneHourAgo)
	if !isInPast(next) {
		t.Fatalf("Unexpected next run time '%s'. Expected to be in the past", next)
	}

	next = f.NextRunTime(halfHourAgo)
	if !dateMatches(next, halfHourFromNow) {
		t.Fatalf("Unexpected next run time. Expected '%s' got '%s'", halfHourFromNow, next)
	}

	next = f.NextRunTime(time.Now())
	if !dateMatches(next, oneHourFromNow) {
		t.Fatalf("Unexpected next run time. Expected '%s' got '%s'", oneHourFromNow, next)
	}
}

func TestEveryDay(t *testing.T) {
	f := freq.Frequency{
		Every:  freq.EveryDay,
		AtTime: midnightLastWeek,
	}

	next := f.NextRunTime(midnightYesterday)
	if !isInPast(next) {
		t.Fatalf("Unexpected next run time '%s'. Expected to be in the past", next)
	}

	next = f.NextRunTime(time.Now())
	if !dateMatches(next, midnightTomorrow) {
		t.Fatalf("Unexpected next run time. Expected '%s' got '%s'", midnightTomorrow, next)
	}
}

func TestEveryWeek(t *testing.T) {
	f := freq.Frequency{
		Every:  freq.EveryWeek,
		AtTime: midnightLastWeek.AddDate(0, 0, -14),
	}

	next := f.NextRunTime(midnightLastWeek.AddDate(0, 0, -7))
	if !isInPast(next) {
		t.Fatalf("Unexpected next run time '%s'. Expected to be in the past", next)
	}

	next = f.NextRunTime(midnightLastWeek)
	if !isInPast(next) {
		t.Fatalf("Unexpected next run time '%s'. Expected to be in the past", next)
	}

	next = f.NextRunTime(time.Now())
	if !dateMatches(next, midnightNextWeek) {
		t.Fatalf("Unexpected next run time. Expected '%s' got '%s'", midnightNextWeek, next)
	}
}

func TestEveryMonth(t *testing.T) {
	firstOfLastMonth := time.Now().AddDate(0, -1, -time.Now().Day()+1)
	firstOfCurrentMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	firstOfNextMonth := time.Now().AddDate(0, 1, -time.Now().Day()+1)
	fmt.Printf("firstOfLastMonth: %s\n", firstOfLastMonth)
	fmt.Printf("firstOfCurrentMonth: %s\n", firstOfCurrentMonth)

	f := freq.Frequency{
		Every:  freq.EveryMonth,
		AtTime: firstOfLastMonth,
	}

	next := f.NextRunTime(firstOfCurrentMonth.AddDate(0, 0, -5))
	if !dateMatches(next, firstOfNextMonth) {
		t.Fatalf("Unexpected next run time. Expected '%s' got '%s'", firstOfNextMonth, next)
	}
}
