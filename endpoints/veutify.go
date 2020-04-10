package endpoints

import (
	"log"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/arunsworld/template/services"
)

// EnableVeutify enables experimental veutify template
func EnableVeutify(srvMux *http.ServeMux, templates []string, ss services.SessionStore) error {
	tmpl := newHTMLFromTemplateFromMinfiedTemplates(templates, "vuetify")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("while processing request %s: %v", r.URL.Path, err)
		}
	})
	srvMux.Handle("/vuetify", gziphandler.GzipHandler(handler))

	return nil
}
