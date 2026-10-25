package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		iso  string
		want string
	}{
		{"2026-10-27T21:30:00", "21.30"},
		{"2026-10-27T23:00:00", "23.00"},
		{"2026-01-02T09:05:00", "09.05"},
	}
	for _, tt := range tests {
		if got := formatTime(tt.iso); got != tt.want {
			t.Errorf("formatTime(%q) = %q, want %q", tt.iso, got, tt.want)
		}
	}
}

func TestFormatTimeInvalidReturnsInput(t *testing.T) {
	raw := "not-a-time"
	if got := formatTime(raw); got != raw {
		t.Errorf("formatTime(%q) = %q, want input back", raw, got)
	}
}

func TestGroupAvailabilities(t *testing.T) {
	slots := []availability{
		{Time: "2026-10-27T17:30:00", Shift: shift{Name: "Lunch"}},
		{Time: "2026-10-27T18:00:00", Shift: shift{Name: "Lunch"}},
		{Time: "2026-10-27T21:30:00", Shift: shift{Name: "Dinner"}},
		{Time: "2026-10-27T22:00:00", Shift: shift{Name: "Lunch"}},
	}
	groups := groupAvailabilities(slots)
	if len(groups) != 2 {
		t.Fatalf("len(groups) = %d, want 2", len(groups))
	}
	if groups[0].Name != "Lunch" {
		t.Errorf("groups[0].Name = %q, want Lunch (first-seen order)", groups[0].Name)
	}
	if len(groups[0].Times) != 3 {
		t.Errorf("groups[0].Times = %v, want 3 entries", groups[0].Times)
	}
	if groups[0].Times[0] != "17.30" {
		t.Errorf("groups[0].Times[0] = %q, want 17.30", groups[0].Times[0])
	}
	if groups[1].Name != "Dinner" || len(groups[1].Times) != 1 {
		t.Errorf("groups[1] = %+v, want single Dinner group", groups[1])
	}
}

func TestGroupAvailabilitiesEmpty(t *testing.T) {
	if groups := groupAvailabilities(nil); groups != nil {
		t.Errorf("groupAvailabilities(nil) = %+v, want nil", groups)
	}
}

func TestRenderPage(t *testing.T) {
	recorder := httptest.NewRecorder()
	renderPage(recorder, pageData{
		VenueName:  "BasBas Kulma",
		BookingURL: "https://book.tableonline.com/mod_0VWSJZ5irGQMNZ",
		Nav: []navItem{
			{Href: "/", Label: "BasBas"},
		},
		Days: []dayArticle{
			{
				Heading: "TIISTAI 27.10",
				Groups: []shiftGroup{
					{Name: "Dinner", Times: []string{"21.30", "22.00"}},
				},
				BookingURL: "https://book.tableonline.com/mod_0VWSJZ5irGQMNZ",
			},
			{
				Heading:    "TORSTAI 29.10",
				Groups:     nil,
				BookingURL: "https://book.tableonline.com/mod_0VWSJZ5irGQMNZ",
			},
			{
				Heading:    "PERJANTAI 30.10",
				Failed:     true,
				BookingURL: "https://book.tableonline.com/mod_0VWSJZ5irGQMNZ",
			},
		},
	})

	body := recorder.Body.String()
	for _, want := range []string{
		"BasBas Kulma p&ouml;yt&auml;&shy;varaukset",
		"TIISTAI 27.10",
		"<h2>Dinner</h2>",
		"<li>21.30</li>",
		"<li>22.00</li>",
		"Ei vapaita aikoja.",
		"Emme voineet hakea vapaita aikoja PERJANTAI 30.10.",
		"https://book.tableonline.com/mod_0VWSJZ5irGQMNZ",
		`<a href="/">BasBas p&ouml;yt&auml;&shy;varaukset</a>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("render output missing %q", want)
		}
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q", got)
	}
}
