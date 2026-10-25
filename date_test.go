package main

import (
	"testing"
	"time"
)

func TestISODate(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{"regular", dateInVenueLocation(t, 2026, 10, 27), "2026-10-27"},
		{"single digits", dateInVenueLocation(t, 2026, 1, 5), "2026-01-05"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isoDate(tt.date); got != tt.want {
				t.Errorf("isoDate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDisplayDate(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{"tuesday", dateInVenueLocation(t, 2026, 10, 27), "TIISTAI 27.10"},
		{"single digits", dateInVenueLocation(t, 2026, 11, 3), "TIISTAI 3.11"},
		{"sunday", dateInVenueLocation(t, 2026, 10, 25), "SUNNUNTAI 25.10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := displayDate(tt.date); got != tt.want {
				t.Errorf("displayDate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAddDays(t *testing.T) {
	start := dateInVenueLocation(t, 2026, 10, 25)
	got := addDays(start, 2)
	if isoDate(got) != "2026-10-27" {
		t.Errorf("addDays() = %q, want 2026-10-27", isoDate(got))
	}
	if got.Hour() != 0 || got.Minute() != 0 {
		t.Errorf("addDays() lost midnight: %v", got)
	}
}

func TestAddDaysAcrossDST(t *testing.T) {
	before := dateInVenueLocation(t, 2026, 3, 27)
	after := addDays(before, 3)
	if isoDate(after) != "2026-03-30" {
		t.Errorf("addDays() across DST = %q, want 2026-03-30", isoDate(after))
	}
	if after.Hour() != 0 {
		t.Errorf("addDays() across DST lost midnight: %v", after)
	}
}

func TestIsWeekdayOpen(t *testing.T) {
	tests := []struct {
		day  time.Time
		want bool
	}{
		{dateInVenueLocation(t, 2026, 10, 25), false},
		{dateInVenueLocation(t, 2026, 10, 26), false},
		{dateInVenueLocation(t, 2026, 10, 27), true},
		{dateInVenueLocation(t, 2026, 10, 28), true},
		{dateInVenueLocation(t, 2026, 10, 29), true},
		{dateInVenueLocation(t, 2026, 10, 30), true},
		{dateInVenueLocation(t, 2026, 10, 31), false},
	}
	for _, tt := range tests {
		if got := isWeekdayOpen(tt.day); got != tt.want {
			t.Errorf("isWeekdayOpen(%s) = %v, want %v", isoDate(tt.day), got, tt.want)
		}
	}
}

func TestMidnightConvertsForeignTimezone(t *testing.T) {
	utc := time.Date(2026, 10, 26, 23, 30, 0, 0, time.UTC)
	got := midnight(utc)
	if isoDate(got) != "2026-10-27" {
		t.Errorf("midnight() = %q, want 2026-10-27 (Helsinki is ahead of UTC)", isoDate(got))
	}
}

func dateInVenueLocation(t *testing.T, year int, month time.Month, day int) time.Time {
	t.Helper()
	return time.Date(year, month, day, 0, 0, 0, 0, venueLocation)
}
