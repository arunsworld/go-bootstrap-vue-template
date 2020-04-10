package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/arunsworld/template/services"
)

// EnableLogin enables the home page
func EnableLogin(srvMux *http.ServeMux, templates []string, ss services.SessionStore, auth services.Auth) error {
	tmpl := newHTMLFromTemplateFromMinfiedTemplates(templates, "login-vuetify")

	loginHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if err := tmpl.Execute(w, nil); err != nil {
				log.Printf("while processing request %s: %v", r.URL.Path, err)
			}
		case "POST":
			processLoginRequest(w, r, ss, auth)
		}
	})
	srvMux.Handle("/login", gziphandler.GzipHandler(loginHandler))

	srvMux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if err := ss.Logout(r, w); err != nil {
			log.Printf("while processing request %s: %v", r.URL.Path, err)
			return
		}
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	})

	return nil
}

func processLoginRequest(w http.ResponseWriter, r *http.Request, ss services.SessionStore, auth services.Auth) {
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
	userID, ok, err := authenticate(auth, username, password)
	if err != nil {
		log.Println(err)
		http.Error(w, "There was a problem on the server side", http.StatusInternalServerError)
		return
	}
	switch {
	case ok:
		if err := ss.StoreUserID(userID, r, w); err != nil {
			log.Println(err)
			http.Error(w, "There was a problem on the server side", http.StatusInternalServerError)
			return
		}
		nextURL, err := ss.NextURL(r, w)
		if err != nil {
			log.Println(err)
			http.Error(w, "There was a problem on the server side", http.StatusInternalServerError)
			return
		}
		resp := struct {
			Next string
		}{
			Next: nextURL,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	default:
		http.Error(w, "credentials incorrect", http.StatusForbidden)
	}
}

func authenticate(auth services.Auth, username, password string) (string, bool, error) {
	ok, err := auth.Authenticate(username, password)
	if err != nil {
		return "", false, err
	}
	return username, ok, nil
}
