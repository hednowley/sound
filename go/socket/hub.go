package socket

import (
	"encoding/json"
	"net/http"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/socket/dto"
)

// IHub maintains the set of active clients and broadcasts messages to the clients.
type IHub interface {
	Notify(notification *dto.Notification)
	SetHandler(method string, handler Handler)

	// Tries to set up a new client and register it with the hub.
	AddClient(w http.ResponseWriter, r *http.Request)
	Run()
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type hub struct {
	ticketer *Ticketer

	// Registered clients.
	clients map[*Client]bool

	// Receieves messages which should be sent out by all clients.
	broadcast chan []byte

	// Receives new clients which want to be registered with this hub.
	register chan *Client

	// Unregisters clients from this hub.
	unregister chan *Client

	// Receives requests forwarded by clients
	incoming chan *incoming

	handlers map[string]Handler
}

type incoming struct {
	client  *Client
	request *dto.Request
}

// NewHub creates a new hub.
func NewHub(ticketer *Ticketer) IHub {
	return &hub{
		ticketer:   ticketer,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		incoming:   make(chan *incoming),
		clients:    make(map[*Client]bool),
		handlers:   make(map[string]Handler),
	}
}

// SetHandler makes sure all messages with the given method are passed to the given handler.
func (h *hub) SetHandler(method string, handler Handler) {
	h.handlers[method] = handler
}

// Run starts the hub.
func (h *hub) Run() {
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

// runHandler takes a request from a client, runs the appropriate handler for
// the request and then send the response back to the client.
func (h *hub) runHandler(incoming *incoming) {

	// Check that this client still registered with the hub.
	registered, ok := h.clients[incoming.client]
	if !registered || !ok {
		return
	}

	handler, ok := h.handlers[incoming.request.Method]
	if ok {
		response := handler(incoming.request)
		if response == nil {
			return
		}

		j, err := json.Marshal(dto.NewResponse(response, incoming.request.ID))
		if err != nil {
			return
		}

		// Check that client is still registered to the hub.
		// It couild have been deregistered since the handler was invoked.
		registered, ok := h.clients[incoming.client]
		if registered && ok {
			incoming.client.send <- j
		}
	}
}

// Notify sends a notification to all clients.
func (h *hub) Notify(notification *dto.Notification) {
	b, err := json.Marshal(notification)
	if err != nil {
		h.broadcast <- []byte("oops")
	}
	h.broadcast <- b
}

// AddClient is a web controller which creates new socket clients.
func (h *hub) AddClient(w http.ResponseWriter, r *http.Request) {

	// Allow all origins
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		conn.Close()
		seelog.Error(err)
		return
	}

	conn.SetReadLimit(maxMessageSize)

	c := NewConnection(conn)
	request, err := c.ReadMessage()
	if err != nil {
		conn.Close()
		seelog.Error(err)
		return
	}

	if request.Method != "handshake" {
		conn.Close()
		seelog.Errorf("Unexpected method: %v", request.Method)
		return
	}

	var ticket string
	err = json.Unmarshal(*request.Params["ticket"], &ticket)
	if err != nil {
		conn.Close()
		seelog.Error("Unexpected handshake params")
		return
	}

	user := h.ticketer.SubmitTicket(ticket)
	if user == nil {
		c.SendMessage(dto.NewErrorResponse("Bad ticket", request.ID))
		conn.Close()
		return
	}

	c.SendMessage(dto.NewResponse(dto.TicketResponse{
		Accepted: true,
	}, request.ID))

	client := &Client{
		hub:  h,
		conn: c,
		send: make(chan []byte, 256),
	}

	h.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
