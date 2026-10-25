package main

import (
	"time"
	_ "time/tzdata"
)

var venueLocation = mustLocation("Europe/Helsinki")

func mustLocation(name string) *time.Location {
	location, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return location
}

func midnight(date time.Time) time.Time {
	date = date.In(venueLocation)
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, venueLocation)
}

func venueToday() time.Time {
	return midnight(time.Now())
}

func addDays(date time.Time, days int) time.Time {
	return midnight(date.AddDate(0, 0, days))
}

func isoDate(date time.Time) string {
	return date.Format("2006-01-02")
}

var weekdays = [...]string{
	"SUNNUNTAI",
	"MAANANTAI",
	"TIISTAI",
	"KESKIVIIKKO",
	"TORSTAI",
	"PERJANTAI",
	"LAUANTAI",
}

func displayDate(date time.Time) string {
	return weekdays[date.Weekday()] + " " + date.Format("2.1")
}

func isWeekdayOpen(date time.Time) bool {
	switch date.Weekday() {
	case time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		return true
	default:
		return false
	}
}
