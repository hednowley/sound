package ws

import (
	"encoding/json"
	"errors"

	"github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"github.com/hednowley/sound/ws/dto"
)

// Connection is a wrapper around a websocket connection.
type Connection struct {
	Inner *websocket.Conn
}

// NewConnection creates a new connection.
func NewConnection(inner *websocket.Conn) *Connection {
	return &Connection{
		Inner: inner,
	}
}

// SendMessage sends the response to the remote client.
func (c *Connection) SendMessage(r *dto.Response) {

	body, err := json.Marshal(r)
	if err != nil {

	}

	c.Inner.WriteMessage(websocket.TextMessage, body)
}

// ReadMessage returns the last message sent from the remote client.
func (c *Connection) ReadMessage() (*dto.Request, error) {
	messageType, payload, err := c.Inner.ReadMessage()
	if err != nil {
		return nil, err
	}

	if messageType != websocket.TextMessage {
		return nil, errors.New("Non-text message received")
	}

	var r dto.Request
	err = json.Unmarshal(payload, &r)
	if err != nil {
		seelog.Errorf("Unexpected request: %v", string(payload))
		return nil, err
	}

	return &r, nil
}

// Close closes the connection.
func (c *Connection) Close() error {
	return c.Inner.Close()
}
