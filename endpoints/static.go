package endpoints

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

// EnableStatic enables serving a static fileserver endpoint
// typical use: EnableStatic("/static/", "static")
func EnableStatic(srvMux *http.ServeMux, path string, dir string) {
	fs := http.FileServer(safeFileSystem{fs: http.Dir(dir)})
	srvMux.Handle(path, http.StripPrefix(path, fs))
}

// safeFileSystem prevents directory listing
type safeFileSystem struct {
	fs http.FileSystem
}

func (sfs safeFileSystem) Open(path string) (http.File, error) {
	f, err := sfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

// EnableFaviconIco serves favicon.ico
func EnableFaviconIco(srvMux *http.ServeMux, staticPath string) error {
	f, err := os.Open(path.Join(staticPath, "favicon.ico"))
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	srvMux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		_, err := w.Write(content)
		if err != nil {
			log.Printf("ERROR: %v", err)
			http.Error(w, "A problem occurred", http.StatusInternalServerError)
			return
		}
	})
	return nil
}

// EnableRobots enables robots.txt
func EnableRobots(srvMux *http.ServeMux) {
	content := []byte(`user-agent: *
allow: /`)
	srvMux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(content)
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
	})
}
