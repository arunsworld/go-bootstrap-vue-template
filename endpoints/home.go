package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/websocket"

	"github.com/arunsworld/template/services"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		origin, ok := r.Header["Origin"]
		if !ok {
			return false
		}
		if strings.HasPrefix(origin[0], "https://bootstrap.iapps365.com") {
			return true
		}
		if strings.HasPrefix(origin[0], "http://localhost") {
			return true
		}
		return false
	}
}

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

	srvMux.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		for i := 0; i < 1000; i++ {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(i+1))); err != nil {
				log.Println(err)
				return
			}
			time.Sleep(time.Second)
		}
	})

	srvMux.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
		log.Println("switch called on:", Theme)
		switch Theme {
		case "bootstrap":
			tmpl = newHTMLFromTemplateFromMinfiedTemplates(templates, "home-vuetify")
			Theme = "vuetify"
		case "vuetify":
			tmpl = newHTMLFromTemplateFromMinfiedTemplates(templates, "home-bootstrap")
			Theme = "bootstrap"
		}
		log.Println("Theme:", Theme)
	})

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
