package main

import (
	"encoding/json"
	"testing"
)

const calendarFixture = `[
	{"date":"2026-10-27","persons":{"2":{"total_slots":464,"available_slots":27,"available_times":5,"waitlist":false,"promotions":[]}}},
	{"date":"2026-10-28","persons":{"2":{"total_slots":480,"available_slots":64,"available_times":9,"waitlist":false,"promotions":[]}}},
	{"date":"2026-11-01","persons":{"2":{"total_slots":0,"available_slots":0,"available_times":0,"waitlist":false,"promotions":[]}}}
]`

const dayFixture = `[
	{"time":"2026-10-27T21:30:00","shift":{"name":"Dinner","category":["dinner"]},"tables":[]},
	{"time":"2026-10-27T22:00:00","shift":{"name":"Dinner","category":["dinner"]},"tables":[]},
	{"time":"2026-10-27T23:00:00","shift":{"name":"Dinner","category":["dinner"]},"tables":[]}
]`

func TestParseCalendar(t *testing.T) {
	var days []calendarDay
	if err := json.Unmarshal([]byte(calendarFixture), &days); err != nil {
		t.Fatal(err)
	}
	if len(days) != 3 {
		t.Fatalf("len(days) = %d, want 3", len(days))
	}
	if days[0].Date != "2026-10-27" {
		t.Errorf("Date = %q", days[0].Date)
	}
	if got := days[0].totalSlots(2); got != 464 {
		t.Errorf("totalSlots(2) = %d, want 464", got)
	}
	if got := days[1].availableTimes(2); got != 9 {
		t.Errorf("availableTimes(2) = %d, want 9", got)
	}
	if got := days[2].totalSlots(2); got != 0 {
		t.Errorf("closed day totalSlots(2) = %d, want 0", got)
	}
	if got := days[0].totalSlots(4); got != 0 {
		t.Errorf("missing persons entry totalSlots(4) = %d, want 0", got)
	}
}

func TestCalendarDateParsing(t *testing.T) {
	var days []calendarDay
	if err := json.Unmarshal([]byte(calendarFixture), &days); err != nil {
		t.Fatal(err)
	}
	date, err := days[0].date()
	if err != nil {
		t.Fatal(err)
	}
	if isoDate(date) != "2026-10-27" {
		t.Errorf("date() = %q", isoDate(date))
	}
}

func TestParseDayTimes(t *testing.T) {
	var slots []availability
	if err := json.Unmarshal([]byte(dayFixture), &slots); err != nil {
		t.Fatal(err)
	}
	if len(slots) != 3 {
		t.Fatalf("len(slots) = %d, want 3", len(slots))
	}
	if slots[0].Time != "2026-10-27T21:30:00" {
		t.Errorf("Time = %q", slots[0].Time)
	}
	if slots[0].Shift.Name != "Dinner" {
		t.Errorf("Shift.Name = %q", slots[0].Shift.Name)
	}
}

func TestParseEmptyDay(t *testing.T) {
	var slots []availability
	if err := json.Unmarshal([]byte(`[]`), &slots); err != nil {
		t.Fatal(err)
	}
	if len(slots) != 0 {
		t.Errorf("len(slots) = %d, want 0", len(slots))
	}
}
