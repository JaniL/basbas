package main

import "testing"

func TestVenueBySlug(t *testing.T) {
	tests := []struct {
		slug string
		want string
		ok   bool
	}{
		{"", "BasBas", true},
		{"kulma", "BasBas Kulma", true},
		{"studio", "BasBas Studio", true},
		{"sirkus", "BasBas Sirkus", true},
		{"unknown", "", false},
	}
	for _, tt := range tests {
		v, ok := venueBySlug(tt.slug)
		if ok != tt.ok {
			t.Errorf("venueBySlug(%q) ok = %v, want %v", tt.slug, ok, tt.ok)
			continue
		}
		if v.name != tt.want {
			t.Errorf("venueBySlug(%q) = %q, want %q", tt.slug, v.name, tt.want)
		}
	}
}

func TestAllVenuesHaveIdentifiers(t *testing.T) {
	for _, v := range venues {
		if v.moduleID == "" || v.restaurantID == "" || v.name == "" {
			t.Errorf("venue %q is missing identifiers: %+v", v.slug, v)
		}
	}
}

func TestBookingURL(t *testing.T) {
	v, _ := venueBySlug("kulma")
	if got := v.bookingURL(); got != "https://book.tableonline.com/mod_0VWSJZ5irGQMNZ" {
		t.Errorf("bookingURL() = %q", got)
	}
}

func TestVenueHref(t *testing.T) {
	kulma, _ := venueBySlug("kulma")
	tests := []struct {
		persons int
		want    string
	}{
		{2, "/kulma"},
		{1, "/kulma?persons=1"},
		{8, "/kulma?persons=8"},
	}
	for _, tt := range tests {
		if got := venueHref(kulma, tt.persons); got != tt.want {
			t.Errorf("venueHref(persons=%d) = %q, want %q", tt.persons, got, tt.want)
		}
	}
}
