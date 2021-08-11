package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var WsChan = make(chan WsPayload)

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
