package main

import "testing"

func TestClampPersons(t *testing.T) {
	tests := []struct {
		raw  string
		want int
	}{
		{"", 2},
		{"2", 2},
		{"1", 1},
		{"8", 8},
		{"0", 1},
		{"9", 8},
		{"100", 8},
		{"-3", 1},
		{"abc", 2},
		{"2.5", 2},
	}
	for _, tt := range tests {
		if got := clampPersons(tt.raw); got != tt.want {
			t.Errorf("clampPersons(%q) = %d, want %d", tt.raw, got, tt.want)
		}
	}
}
