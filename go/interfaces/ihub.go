package interfaces

import (
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/ws/dto"
	"net/http"
)

// WsHandler listens for particular websocket messages.
type WsHandler = func(*dto.Request) interface{}

// Ticketer manages tickets for negotiating websocket connections.
type Ticketer interface {
	SubmitTicket(key string) *config.User
}

// Hub manages a selection of clients who can send and receive messages.
type Hub interface {
	Notify(notification *dto.Notification)
	SetHandler(method string, handler WsHandler)
	AddClient(ticketer Ticketer, dal DAL, w http.ResponseWriter, r *http.Request)
	Run()
}
