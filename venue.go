package main

import "strconv"

type venue struct {
	slug         string
	name         string
	moduleID     string
	restaurantID string
}

func (v venue) bookingURL() string {
	return "https://book.tableonline.com/" + v.moduleID
}

var venues = []venue{
	{slug: "", name: "BasBas", moduleID: "mod_0VWSE3bovStZfY", restaurantID: "res_0VWSAUspRta7lU"},
	{slug: "kulma", name: "BasBas Kulma", moduleID: "mod_0VWSJZ5irGQMNZ", restaurantID: "res_0VWSIrekOaBRlH"},
	{slug: "studio", name: "BasBas Studio", moduleID: "mod_0VWSJXzImCZp8l", restaurantID: "res_0VWSJ2xpOSZSvd"},
	{slug: "sirkus", name: "BasBas Sirkus", moduleID: "mod_0VXMjuwnHSmomj", restaurantID: "res_0VXJGaf4He0LPF"},
}

func venueBySlug(slug string) (venue, bool) {
	for _, v := range venues {
		if v.slug == slug {
			return v, true
		}
	}
	return venue{}, false
}

func venueHref(v venue, persons int) string {
	href := "/" + v.slug
	if persons != defaultPersons {
		href += "?persons=" + strconv.Itoa(persons)
	}
	return href
}
