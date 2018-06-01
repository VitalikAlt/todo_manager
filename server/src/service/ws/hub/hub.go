package hub

import (
	"log"

	"github.com/gorilla/websocket"
)

// NewHub - создание новой шины соединений
func NewHub() *Hub {
	hub := &Hub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		clientsByID: make(map[int][]*Client),
	}

	go hub.run()

	return hub
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Clients grouped by id.
	clientsByID map[int][]*Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewClient - создание нового клиента внутри шины
func (h *Hub) NewClient(conn *websocket.Conn) *Client {
	client := &Client{hub: h, conn: conn}
	h.register <- client

	return client
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Printf("[INFO] client hub: client register: %#v", client.conn.RemoteAddr().String())
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Printf("[INFO] client hub: client unregister: %#v", client.conn.RemoteAddr().String())
				delete(h.clients, client)
			}
		}
	}
}
