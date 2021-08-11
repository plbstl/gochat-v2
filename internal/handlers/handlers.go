package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/CloudyKit/jet/v6"
)

var (
	mu      sync.Mutex
	clients = make(map[wsConn]string)
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(),
)

// broadcast sends a response to all connected clients.
func broadcast(r WsResponse) {
	mu.Lock()
	for client := range clients {
		err := client.WriteJSON(r)
		if err != nil {
			log.Println("client.WriteJSON=", err)
			client.Close()
			delete(clients, client)
		}
	}
	mu.Unlock()
}

// renderPage .
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

// users returns a list of connected usernames.
func users() []string {
	usernames := make([]string, 0, len(clients))
	mu.Lock()
	for _, username := range clients {
		usernames = append(usernames, username)
	}
	mu.Unlock()
	return usernames
}
