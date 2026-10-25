package main

import (
	"html/template"
	"net/http"
	"time"
)

type pageData struct {
	VenueName  string
	BookingURL string
	Nav        []navItem
	Days       []dayArticle
	StartedAt  string
}

type navItem struct {
	Href  string
	Label string
}

type dayArticle struct {
	Heading    string
	Groups     []shiftGroup
	Failed     bool
	BookingURL string
}

type shiftGroup struct {
	Name  string
	Times []string
}

const pageHTML = `<!DOCTYPE html>
<html>
<head>
<!-- Server started: {{.StartedAt}}-->
<title>BasBas</title>
<meta name="apple-mobile-web-app-capable" content="yes"/>
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
<link rel="apple-touch-icon" href="menu.png"/>
<link id="page_favicon" href="/favicon.ico" rel="icon" type="image/x-icon"/>
<link rel="stylesheet" type="text/css" href="public/styles.css"/>
</head>
<body>
<article>
<header>
<h1>{{.VenueName}} p&ouml;yt&auml;&shy;varaukset</h1>
<p class="subtitle"><a href="{{.BookingURL}}" target="_blank" rel="noreferrer">Varaa t&auml;&auml;lt&auml;</a></p>
{{range .Nav}}<p class="subtitle"><a href="{{.Href}}">{{.Label}} p&ouml;yt&auml;&shy;varaukset</a></p>
{{end}}</header>
</article>
{{range .Days}}{{template "day" .}}
{{end}}<article><a href="https://github.com/eeroan/basbas" target="_blank" rel="noreferrer">L&auml;hdekoodi</a></article>
</body>
</html>
{{define "day"}}<article>
<h1>{{.Heading}}</h1>
{{if .Failed}}<section><p class="alert">Emme voineet hakea vapaita aikoja {{.Heading}}. <a href="{{.BookingURL}}" target="_blank" rel="noreferrer">Varaa suoraan TableOnline-palvelusta</a>.</p></section>
{{else if not .Groups}}<section><p>Ei vapaita aikoja.</p></section>
{{else}}<section>{{range .Groups}}{{if .Name}}<h2>{{.Name}}</h2>{{end}}<ul>{{range .Times}}<li>{{.}}</li>{{end}}</ul>{{end}}</section>{{end}}
</article>{{end}}
`

var pageTemplate = template.Must(template.New("page").Parse(pageHTML))

func groupAvailabilities(slots []availability) []shiftGroup {
	var groups []shiftGroup
	indexByShift := map[string]int{}
	for _, slot := range slots {
		shift := slot.Shift.Name
		index, ok := indexByShift[shift]
		if !ok {
			groups = append(groups, shiftGroup{Name: shift})
			index = len(groups) - 1
			indexByShift[shift] = index
		}
		groups[index].Times = append(groups[index].Times, formatTime(slot.Time))
	}
	return groups
}

func formatTime(iso string) string {
	parsed, err := time.Parse("2006-01-02T15:04:05", iso)
	if err != nil {
		return iso
	}
	return parsed.Format("15.04")
}

func renderPage(w http.ResponseWriter, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := pageTemplate.Execute(w, data); err != nil {
		return
	}
}
