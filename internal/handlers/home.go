package handlers

import (
	"log"
	"net/http"
)

// Home handles the home route.
func Home(w http.ResponseWriter, r *http.Request) {
	if err := renderPage(w, "home.jet", nil); err != nil {
		log.Println("Home.renderPage:", err)
	}
}
