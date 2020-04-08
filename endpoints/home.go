package endpoints

import (
	"log"
	"net/http"

	"github.com/arunsworld/template/services"
)

// EnableHome enables the home page
func EnableHome(srvMux *http.ServeMux, templates []string, ss services.SessionStore) error {
	tmpl := newHTMLFromTemplateFromMinfiedTemplates(templates, "home")

	type Home struct {
		NavBar NavBar
	}

	srvMux.HandleFunc("/", ss.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		u, err := ss.UserID(r)
		if err != nil {
			log.Println(err)
		}
		home := Home{
			NavBar: NavBar{
				NavItems: []NavItem{
					NavItem{Title: "Link A", Link: "/linkA", IsActive: true},
					NavItem{Title: "Link B", Link: "/linkB", IsActive: false},
				},
				User: u,
			},
		}
		if err := tmpl.Execute(w, home); err != nil {
			log.Printf("while processing request %s: %v", r.URL.Path, err)
		}
	}))

	srvMux.HandleFunc("/linkB", func(w http.ResponseWriter, r *http.Request) {
		u, err := ss.UserID(r)
		if err != nil {
			log.Println(err)
		}
		home := Home{
			NavBar: NavBar{
				NavItems: []NavItem{
					NavItem{Title: "Link A", Link: "/linkA", IsActive: false},
					NavItem{Title: "Link B", Link: "/linkB", IsActive: true},
				},
				User: u,
			},
		}
		if err := tmpl.Execute(w, home); err != nil {
			log.Printf("while processing request %s: %v", r.URL.Path, err)
		}
	})

	return nil
}
