package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrNoSuchUser is a sentinel value to indicate that the user doesn't exist in the store
var ErrNoSuchUser = errors.New("No Such User")

// AuthStoreReader offers functionality to read salted password from store
type AuthStoreReader interface {
	ReadPassword(user string) (string, error)
}

// AuthStoreWriter offers functionality to write salted password for a user
type AuthStoreWriter interface {
	WritePassword(user, password string) error
	UpdatePassword(user, password string) error
}

// AuthStore provides Reader/Writer functionality for Auth Stores
type AuthStore interface {
	AuthStoreReader
	AuthStoreWriter
}

// Auth provides Auth Service using the backing store injected
type Auth struct {
	Store AuthStore
}

// Authenticate authenticates the given user and password against the store
func (auth Auth) Authenticate(user, password string) (bool, error) {
	hashedPwd, err := auth.Store.ReadPassword(user)
	switch {
	case err == ErrNoSuchUser:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("problem reading password from store during authenticate: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password)); err != nil {
		switch {
		case err == bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		case err != nil:
			return false, fmt.Errorf("problem evaluating password: %v", err)
		}
	}
	return true, nil
}

// Register registers the given user and password into the store using salted hash
func (auth Auth) Register(user, password string) error {
	_, err := auth.Store.ReadPassword(user)
	if err != ErrNoSuchUser {
		return fmt.Errorf("user already exists")
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("problem hashing password during registration: %v", err)
	}
	if err := auth.Store.WritePassword(user, string(hashedPwd)); err != nil {
		return fmt.Errorf("problem writing password during registration: %v", err)
	}
	return nil
}

// FileAuthStore is an AuthStore backed by a file
type FileAuthStore struct {
	mu      sync.Mutex
	fname   string
	isDirty bool
	creds   map[string]string
	err     error
}

// NewFileAuthStore gives us a new FileAuthStore
func NewFileAuthStore(fname string) (*FileAuthStore, error) {
	creds, err := readCredsOnLoad(fname)
	if err != nil {
		return nil, err
	}
	result := &FileAuthStore{
		fname: fname,
		creds: creds,
	}
	go result.savePeriodically(time.Second)
	return result, nil
}

func readCredsOnLoad(fname string) (map[string]string, error) {
	_, err := os.Stat(fname)
	if err != nil {
		return make(map[string]string), nil
	}
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	csvr := csv.NewReader(f)
ReadLoop:
	for {
		rec, err := csvr.Read()
		switch {
		case err == io.EOF:
			break ReadLoop
		case err != nil:
			return nil, err
		}
		result[rec[0]] = rec[1]
	}
	return result, nil
}

func (store *FileAuthStore) savePeriodically(frequencyDelay time.Duration) {
	for {
		time.Sleep(frequencyDelay)
		if !store.isDirty {
			continue
		}
		log.Println("store is dirty going to save...")
		f, err := ioutil.TempFile("", "pwd")
		if err != nil {
			store.err = err
			log.Println(err)
			return
		}
		csvw := csv.NewWriter(f)
		store.mu.Lock()
		for user, pwd := range store.creds {
			if err := csvw.Write([]string{user, pwd}); err != nil {
				store.err = err
				store.mu.Unlock()
				log.Println(err)
				return
			}
		}
		csvw.Flush()
		f.Close()
		if err := os.Rename(f.Name(), store.fname); err != nil {
			store.err = err
			store.mu.Unlock()
			log.Println(err)
			return
		}
		store.isDirty = false
		store.mu.Unlock()
	}
}

// ReadPassword reads password for the user from the file
func (store *FileAuthStore) ReadPassword(user string) (string, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	pwd, ok := store.creds[user]
	if !ok {
		return "", ErrNoSuchUser
	}

	return pwd, nil
}

// WritePassword writes password to file
func (store *FileAuthStore) WritePassword(user, password string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.creds[user] = password
	store.isDirty = true

	return nil
}

// UpdatePassword updates password in file
func (store *FileAuthStore) UpdatePassword(user, password string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.creds[user] = password
	store.isDirty = true

	return nil
}
