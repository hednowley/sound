package api

import (
	"net/http"
	"os"
	"path/filepath"
)

// ServeSinglePageApp serves static files if they exist, otherwise it serves a single HTML page.
func ServeSinglePageApp(dir string, html string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = os.Stat(filepath.Join(dir, path))
		if os.IsNotExist(err) {
			// Serve the index file
			http.ServeFile(w, r, html)
			return
		} else if err != nil {
			// Some other error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Requested file must exist so serve it
		http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
	}
}
