package ws

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/ws/dto"
	"github.com/hednowley/sound/ws/handlers"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Receives requests forward by clients
	incoming chan *Incoming

	handlers map[string]handlers.WsHandler
}

type Incoming struct {
	client  *Client
	request *dto.Request
}

func NewHub(dal *dal.DAL) *Hub {

	allHandlers := make(map[string]handlers.WsHandler)
	allHandlers["getArtists"] = handlers.MakeGetArtistsHandler(dal)
	allHandlers["startScan"] = handlers.MakeStartScanHandler(dal)

	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		incoming:   make(chan *Incoming),
		clients:    make(map[*Client]bool),
		handlers:   allHandlers,
	}
}

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
			handler, ok := h.handlers[incoming.request.Method]
			if ok {
				go handler(incoming.request)
			}
		}
	}
}
