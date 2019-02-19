package ws

import (
	"encoding/json"
	"errors"

	"github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"github.com/hednowley/sound/ws/dto"
)

type Connection struct {
	inner *websocket.Conn
}

func NewConnection(inner *websocket.Conn) *Connection {
	return &Connection{
		inner: inner,
	}
}

func (c *Connection) SendMessage(r *dto.Response) {

	body, err := json.Marshal(r)
	if err != nil {

	}

	c.inner.WriteMessage(websocket.TextMessage, body)
}

func (c *Connection) ReadMessage() (*dto.Request, error) {
	messageType, payload, err := c.inner.ReadMessage()
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

func (c *Connection) Close() error {
	return c.inner.Close()
}
