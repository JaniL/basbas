package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	apiBaseURL       = "https://api.tableonline.com"
	requestTimeout   = 8 * time.Second
	calendarEndpoint = apiBaseURL + "/restaurant/%s/availability/calendar?start_date=%s&end_date=%s&persons=%d&module_id=%s"
	dayEndpoint      = apiBaseURL + "/restaurant/%s/availability?date=%s&persons=%d&module_id=%s&locale=fi"
)

var apiClient = &http.Client{Timeout: requestTimeout}

type calendarDay struct {
	Date    string                  `json:"date"`
	Persons map[string]calendarSlot `json:"persons"`
}

type calendarSlot struct {
	TotalSlots     int  `json:"total_slots"`
	AvailableSlots int  `json:"available_slots"`
	AvailableTimes int  `json:"available_times"`
	Waitlist       bool `json:"waitlist"`
}

type shift struct {
	Name string `json:"name"`
}

type availability struct {
	Time  string `json:"time"`
	Shift shift  `json:"shift"`
}

func (d calendarDay) slot(persons int) calendarSlot {
	return d.Persons[strconv.Itoa(persons)]
}

func (d calendarDay) totalSlots(persons int) int {
	return d.slot(persons).TotalSlots
}

func (d calendarDay) availableTimes(persons int) int {
	return d.slot(persons).AvailableTimes
}

func (d calendarDay) date() (time.Time, error) {
	return time.ParseInLocation("2006-01-02", d.Date, venueLocation)
}

func calendarURL(v venue, start, end time.Time, persons int) string {
	return fmt.Sprintf(calendarEndpoint, v.restaurantID, isoDate(start), isoDate(end), persons, v.moduleID)
}

func dayURL(v venue, date time.Time, persons int) string {
	return fmt.Sprintf(dayEndpoint, v.restaurantID, isoDate(date), persons, v.moduleID)
}

func fetchCalendar(v venue, start, end time.Time, persons int) ([]calendarDay, error) {
	resp, err := apiClient.Get(calendarURL(v, start, end, persons))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("calendar request failed with status %d", resp.StatusCode)
	}
	var days []calendarDay
	if err := json.NewDecoder(resp.Body).Decode(&days); err != nil {
		return nil, err
	}
	return days, nil
}

func fetchDayTimes(v venue, date time.Time, persons int) ([]availability, error) {
	resp, err := apiClient.Get(dayURL(v, date, persons))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("day request failed with status %d", resp.StatusCode)
	}
	var slots []availability
	if err := json.NewDecoder(resp.Body).Decode(&slots); err != nil {
		return nil, err
	}
	return slots, nil
}
