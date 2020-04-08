package endpoints

import (
	"log"
	"net/http"

	"github.com/arunsworld/template/services"
)

// EnableRegister enables the registration endpoint
func EnableRegister(srvMux *http.ServeMux, templates []string, ss services.SessionStore, auth services.Auth) error {
	tmpl := newHTMLFromTemplateFromMinfiedTemplates(templates, "register")

	srvMux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if err := tmpl.Execute(w, nil); err != nil {
				log.Printf("while processing request %s: %v", r.URL.Path, err)
			}
		case "POST":
			processRegistrationRequest(w, r, ss, auth)
		}
	})

	return nil
}

func processRegistrationRequest(w http.ResponseWriter, r *http.Request, ss services.SessionStore, auth services.Auth) {
	if err := r.ParseForm(); err != nil {
		log.Printf("while processing request %s: %v", r.URL.Path, err)
		http.Error(w, "internal server issue", http.StatusInternalServerError)
		return
	}
	usernames, ok := r.Form["username"]
	if !ok {
		log.Printf("while processing request %s: username not found", r.URL.Path)
		http.Error(w, "username not found", http.StatusBadRequest)
		return
	}
	username := usernames[0]
	passwords, ok := r.Form["password"]
	if !ok {
		log.Printf("while processing request %s: password not found", r.URL.Path)
		http.Error(w, "password not found", http.StatusBadRequest)
		return
	}
	password := passwords[0]
	if err := auth.Register(username, password); err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}
		log.Printf("while processing request %s: unable to register: %v", r.URL.Path, err)
		http.Error(w, "unable to register", http.StatusBadRequest)
		return
	}
	ss.StoreUserID(username, r, w)
}
