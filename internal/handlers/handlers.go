package handlers

import (
	"log"
	"net/http"
	"sync"

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
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"messageType"`
	ConnectedUsers []string `json:"connectedUsers"`
}

// WsPayload defines the payload
// received from client.
type WsPayload struct {
	Action   string `json:"action"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Conn     wsConn `json:"-"`
}

type wsConn struct {
	*websocket.Conn
}

var (
	mu      sync.Mutex
	WsChan  = make(chan WsPayload)
	clients = make(map[wsConn]string)
)

// Websocket handles a websocket connection.
func Websocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader.Upgrade=", err)
	}

	conn := wsConn{ws}
	go listenForWsConn(&conn)
}

// listenForWsConn indefinitely listens for new messages
// coming through the WebSocket connection. It processes
// and sends the messages as payload to WsChan.
func listenForWsConn(conn *wsConn) {
	defer func() {
		conn.Close()
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			return
		}
		payload.Conn = *conn
		WsChan <- payload
	}
}

// ListenOnWsChan indefintely listens for new payload on
// WsChan and processes it according to payload.Action.
func ListenOnWsChan() {
	var response WsResponse
	for {
		p := <-WsChan
		switch p.Action {
		case "change_username":
			mu.Lock()
			clients[p.Conn] = p.Username
			log.Println("Size of Connection Pool:", len(clients))
			mu.Unlock()
			response.Action = "list_users"
			response.ConnectedUsers = users()
			broadcast(response)

		case "disconnect_client":
			mu.Lock()
			delete(clients, p.Conn)
			log.Println("Size of Connection Pool:", len(clients))
			mu.Unlock()
			response.Action = "list_users"
			response.ConnectedUsers = users()
			broadcast(response)

		default:
			response.Action = "unsupported"
			p.Conn.WriteJSON(response)
		}
	}
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

// broadcast sends a response to all connected clients..
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
