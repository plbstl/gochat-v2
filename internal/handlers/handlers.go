package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(),
)

// Home handles the home route.
func Home(w http.ResponseWriter, r *http.Request) {
	if err := renderPage(w, "home.jet", nil); err != nil {
		log.Println("renderPage:", err)
	}
}

func renderPage(w http.ResponseWriter, tmpl string, vars jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		return err
	}
	if err = view.Execute(w, vars, nil); err != nil {
		return err
	}
	return nil
}
