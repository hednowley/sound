package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/ws/dto"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 100 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *Connection

	// Buffered channel of outbound messages.
	send chan []byte

	dal interfaces.DAL
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		request, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.hub.incoming <- &Incoming{
			client:  c,
			request: request,
		}

		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {

	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.Inner.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.Inner.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.Inner.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}

}

// NewClient tries to set up a new client and register it with the hub.
func NewClient(hub *Hub, ticketer *Ticketer, dal interfaces.DAL, w http.ResponseWriter, r *http.Request) {

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

	user := ticketer.SubmitTicket(ticket)
	if user == nil {
		c.SendMessage(dto.NewErrorResponse("Bad ticket", request.ID))
		conn.Close()
		return
	}

	c.SendMessage(dto.NewResponse(dto.TicketResponse{
		Accepted: true,
	}, request.ID))

	client := &Client{
		hub:  hub,
		conn: c,
		send: make(chan []byte, 256),
		dal:  dal,
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
