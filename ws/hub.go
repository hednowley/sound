package ws

import (
	"encoding/json"

	"github.com/hednowley/sound/ws/dto"
)

// WsHandler listens for particular websocket messages.
type WsHandler = func(*dto.Request) interface{}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Receieves messages which should be sent out by all clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Receives requests forwarded by clients
	incoming chan *Incoming

	handlers map[string]WsHandler
}

type Incoming struct {
	client  *Client
	request *dto.Request
}

// NewHub creates a new hub.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		incoming:   make(chan *Incoming),
		clients:    make(map[*Client]bool),
		handlers:   make(map[string]WsHandler),
	}
}

func (h *Hub) SetHandler(method string, handler WsHandler) {
	h.handlers[method] = handler
}

// Run starts the hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		case incoming := <-h.incoming:
			go h.runHandler(incoming)
		}
	}
}

func (h *Hub) runHandler(incoming *Incoming) {
	handler, ok := h.handlers[incoming.request.Method]
	if ok {
		handler(incoming.request)
	}
}

// Notify sends a notification to all clients.
func (h *Hub) Notify(notification *dto.Notification) {
	b, err := json.Marshal(notification)
	if err != nil {
		h.broadcast <- []byte("oops")
	}
	h.broadcast <- b
}
