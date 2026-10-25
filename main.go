package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultPort    = "8080"
	dayWindow      = 33
	cacheTTL       = 60 * time.Second
	minPersons     = 1
	maxPersons     = 8
	defaultPersons = 2
)

var (
	calendarCache     = newTTLCache[[]calendarDay](cacheTTL)
	availabilityCache = newTTLCache[[]availability](cacheTTL)
)

func clampPersons(raw string) int {
	n, err := strconv.Atoi(raw)
	if err != nil {
		return defaultPersons
	}
	if n < minPersons {
		return minPersons
	}
	if n > maxPersons {
		return maxPersons
	}
	return n
}

func handleRequest(startedAt string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		uri := r.URL.Path
		if uri == "/public/styles.css" || uri == "/favicon.ico" || uri == "/menu.png" ||
			strings.HasPrefix(uri, "/public/") {
			serveStatic(w, uri)
			return
		}

		query := r.URL.Query()
		if uri == "/" && query.Has("kulma") {
			redirect := "/kulma"
			if persons := query.Get("persons"); persons != "" {
				redirect += "?persons=" + persons
			}
			http.Redirect(w, r, redirect, http.StatusMovedPermanently)
			return
		}

		v, ok := venueBySlug(strings.Trim(uri, "/"))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		persons := clampPersons(query.Get("persons"))
		renderPage(w, pageData{
			VenueName:  v.name,
			BookingURL: v.bookingURL(),
			Nav:        navItems(persons, v),
			Days:       buildDays(v, persons),
			StartedAt:  startedAt,
		})
	}
}

func navItems(persons int, current venue) []navItem {
	var items []navItem
	for _, v := range venues {
		if v.slug == current.slug {
			continue
		}
		items = append(items, navItem{
			Href:  venueHref(v, persons),
			Label: v.name,
		})
	}
	return items
}

func buildDays(v venue, persons int) []dayArticle {
	start := venueToday()
	end := addDays(start, dayWindow-1)

	days, err := fetchCalendarCached(v, start, end, persons)
	if err != nil {
		log.Printf("[calendar] fetch failed for %s persons=%d: %v", v.slug, persons, err)
		return weekdayFallbackDays(v, start, persons)
	}

	var dates []time.Time
	for _, day := range days {
		if day.totalSlots(persons) <= 0 {
			continue
		}
		date, err := day.date()
		if err != nil {
			continue
		}
		dates = append(dates, date)
	}
	return dayArticles(v, dates, persons)
}

func weekdayFallbackDays(v venue, start time.Time, persons int) []dayArticle {
	var dates []time.Time
	for i := 0; i < dayWindow; i++ {
		date := addDays(start, i)
		if isWeekdayOpen(date) {
			dates = append(dates, date)
		}
	}
	return dayArticles(v, dates, persons)
}

func dayArticles(v venue, dates []time.Time, persons int) []dayArticle {
	articles := make([]dayArticle, len(dates))
	var wg sync.WaitGroup
	for i, date := range dates {
		articles[i] = dayArticle{
			Heading:    displayDate(date),
			BookingURL: v.bookingURL(),
		}
		wg.Add(1)
		go func(i int, date time.Time) {
			defer wg.Done()
			articles[i].Groups, articles[i].Failed = fetchDayGroups(v, date, persons)
		}(i, date)
	}
	wg.Wait()

	failed := 0
	for _, article := range articles {
		if article.Failed {
			failed++
		}
	}
	if failed > 0 {
		log.Printf("[availability] %d days failed for persons=%d venue=%s", failed, persons, v.slug)
	}
	return articles
}

func fetchCalendarCached(v venue, start, end time.Time, persons int) ([]calendarDay, error) {
	key := calendarURL(v, start, end, persons)
	if days, ok := calendarCache.get(key); ok {
		return days, nil
	}
	days, err := fetchCalendar(v, start, end, persons)
	if err != nil {
		return nil, err
	}
	calendarCache.set(key, days)
	return days, nil
}

func fetchDayGroups(v venue, date time.Time, persons int) ([]shiftGroup, bool) {
	key := dayURL(v, date, persons)
	if slots, ok := availabilityCache.get(key); ok {
		return groupAvailabilities(slots), false
	}
	slots, err := fetchDayTimes(v, date, persons)
	if err != nil {
		log.Printf("[availability] day %s failed for %s persons=%d: %v", isoDate(date), v.slug, persons, err)
		return nil, true
	}
	availabilityCache.set(key, slots)
	return groupAvailabilities(slots), false
}

func main() {
	started := time.Now()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	startedAt := started.Format(time.RFC1123)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server started in http://localhost:%s, took %s", port, time.Since(started))

	server := &http.Server{Handler: handleRequest(startedAt)}
	log.Fatal(server.Serve(listener))
}
