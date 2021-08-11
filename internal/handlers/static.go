package handlers

import "net/http"

// StaticFiles serves files from provided folder.
func StaticFiles(folder string) http.Handler {
	return http.StripPrefix("/static", http.FileServer(http.Dir(folder)))
}
