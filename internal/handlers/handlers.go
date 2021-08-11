package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.InDevelopmentMode(),
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Home handles the home route.
func Home(w http.ResponseWriter, r *http.Request) {
	if err := renderPage(w, "home.jet", nil); err != nil {
		log.Println("renderPage:", err)
	}
}

// WsResponse defines the response
// sent back from websocket.
type WsResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// Websocket handles a websocket connection.
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader.Upgrade:", err)
	}

	log.Println("Connected Client")
	response := WsResponse{
		Message: "Connected to server!",
	}

	err = conn.WriteJSON(response)
	if err != nil {
		log.Println("conn.WriteJSON:", err)
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
