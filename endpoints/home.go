package endpoints

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NYTimes/gziphandler"

	"github.com/arunsworld/template/services"
)

// EnableHome enables the home page
func EnableHome(srvMux *http.ServeMux, templates []string, ss services.SessionStore) error {
	tmpl := newHTMLFromTemplateFromMinfiedTemplates(templates, fmt.Sprintf("home-%s", Theme))

	type Home struct {
		NavBar NavBar
		User   string
	}
	homeNavBar := newNavBarFor("Home")
	dashboardNavBar := newNavBarFor("Dashboard")

	srvMux.Handle("/", gziphandler.GzipHandler(ss.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		u, err := ss.UserID(r)
		if err != nil {
			log.Println(err)
		}
		home := Home{
			NavBar: homeNavBar,
			User:   u,
		}
		if err := tmpl.Execute(w, home); err != nil {
			log.Printf("while processing request %s: %v", r.URL.Path, err)
		}
	})))

	srvMux.Handle("/dashboard", gziphandler.GzipHandler(ss.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		u, err := ss.UserID(r)
		if err != nil {
			log.Println(err)
		}
		home := Home{
			NavBar: dashboardNavBar,
			User:   u,
		}
		if err := tmpl.Execute(w, home); err != nil {
			log.Printf("while processing request %s: %v", r.URL.Path, err)
		}
	})))

	return nil
}

func newNavBarFor(page string) NavBar {
	result := NavBar{
		NavItems: []NavItem{
			NavItem{Title: "Home", Link: "/", Icon: "mdi-home"},
			NavItem{Title: "Dashboard", Link: "/dashboard", Icon: "mdi-view-dashboard"},
			NavItem{Title: "Settings", Link: "/settings", Icon: "mdi-settings"},
		},
	}
	activeItem := 0
	for i, item := range result.NavItems {
		if item.Title == page {
			activeItem = i
			result.NavItems[i].IsActive = true
		}
	}
	result.ActiveItem = activeItem
	return result
}
