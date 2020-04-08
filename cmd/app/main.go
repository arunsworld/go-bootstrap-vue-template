package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/sessions"

	"github.com/arunsworld/template/endpoints"
	"github.com/arunsworld/template/services"
	"github.com/tdewolff/minify/v2"
)

func main() {
	webDir := flag.String("web", "web", "directory where web assets are stored")
	port := flag.Int("port", 0, "port to run server on")
	flag.Parse()
	validate(*webDir)
	if *port == 0 {
		*port = 1000
	}

	staticDir := path.Join(*webDir, "static")
	endpoints.EnableStatic("/static/", staticDir)
	if err := endpoints.EnableFaviconIco(staticDir); err != nil {
		log.Fatal(err)
	}

	ss := sessionStore()
	templates, err := endpoints.MinifiedTemplates(path.Join(*webDir, "templates"))
	if err != nil {
		log.Fatal(err)
	}

	authStore, err := services.NewFileAuthStore("pwd.csv")
	if err != nil {
		log.Fatal(err)
	}
	auth := services.Auth{
		Store: authStore,
	}
	if err := endpoints.EnableLogin(templates, ss, auth); err != nil {
		log.Fatal(err)
	}
	if err := endpoints.EnableHome(templates, ss); err != nil {
		log.Fatal(err)
	}
	if err := endpoints.EnableRegister(templates, ss, auth); err != nil {
		log.Fatal(err)
	}

	services.ServerStart(services.ServerConfig{
		Addr: fmt.Sprintf(":%d", *port),
	})
}

func sessionStore() services.SessionStore {
	sessionKeyEnc := os.Getenv("SESSION_KEY")
	if sessionKeyEnc == "" {
		log.Fatal("cannot proceed without SESSION_KEY")
	}
	sessionKey, err := hex.DecodeString(sessionKeyEnc)
	if err != nil {
		log.Fatal(err)
	}
	return services.NewSessionStore(sessions.NewCookieStore(
		sessionKey,
	), "/login")
}

func minifierMiddleware(minifier *minify.M) services.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		mHandler := minifier.Middleware(next)
		return func(w http.ResponseWriter, r *http.Request) {
			mHandler.ServeHTTP(w, r)
		}
	}
}

func validate(webDir string) {
	if webDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	info, err := os.Stat(webDir)
	if err != nil {
		log.Printf("%s does not exist. cannot start HTTP server.", webDir)
		os.Exit(1)
	}

	if !info.IsDir() {
		log.Printf("%s is not a directory. cannot start HTTP server.", webDir)
		os.Exit(1)
	}
}
