package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// SessionStore is a store for managing sessions and providing required middleware services
type SessionStore struct {
	Store          sessions.Store
	AuthMiddleware Middleware
}

// Middleware is a middleware at the HandlerFunc level
type Middleware func(http.HandlerFunc) http.HandlerFunc

// NewSessionStore returns a new session store
func NewSessionStore(store sessions.Store, loginURL string) SessionStore {
	result := SessionStore{
		Store: store,
	}
	result.AuthMiddleware = result.authMiddleware(loginURL)
	return result
}

// authMiddleware returns a middelware that ensures session is authenticated
func (ss SessionStore) authMiddleware(loginURL string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session, err := ss.Store.Get(r, "auth")
			if err != nil {
				log.Println(err)
				http.Error(w, "There was a problem on the server side", http.StatusInternalServerError)
				return
			}
			userID := session.Values["userID"]
			switch userID == nil {
			case true:
				nextURL := r.URL.Path
				if nextURL != loginURL {
					log.Printf("found nextURL: %s", nextURL)
					session.AddFlash(nextURL)
					if err := session.Save(r, w); err != nil {
						log.Println("unable to save session")
					}
				}
				http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
			default:
				next(w, r)
			}
		}
	}
}

// StoreUserID stores the userID with the cookie associated with request
// this should be called once the credentials are validated
func (ss SessionStore) StoreUserID(userID string, r *http.Request, w http.ResponseWriter) error {
	session, err := ss.Store.Get(r, "auth")
	if err != nil {
		return err
	}
	session.Values["userID"] = userID
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

// UserID retreives the userID from the session
func (ss SessionStore) UserID(r *http.Request) (string, error) {
	session, err := ss.Store.Get(r, "auth")
	if err != nil {
		return "", err
	}
	userID := session.Values["userID"]
	switch userID == nil {
	case true:
		return "", fmt.Errorf("UserID not found")
	default:
		return userID.(string), nil
	}
}

// NextURL retreives the nextURL from the session flash
func (ss SessionStore) NextURL(r *http.Request, w http.ResponseWriter) (string, error) {
	session, err := ss.Store.Get(r, "auth")
	if err != nil {
		return "", err
	}
	if flashes := session.Flashes(); len(flashes) > 0 {
		session.Save(r, w)
		return flashes[0].(string), nil
	}
	return "/", nil
}

// Logout removes the userID from the session
func (ss SessionStore) Logout(r *http.Request, w http.ResponseWriter) error {
	session, err := ss.Store.Get(r, "auth")
	if err != nil {
		return err
	}
	delete(session.Values, "userID")
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}
