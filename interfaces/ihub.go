package interfaces

import (
	"net/http"

	"github.com/hednowley/sound/ws/dto"
)

// WsHandler listens for particular websocket messages.
type WsHandler = func(*dto.Request) interface{}

type Hub interface {
	Notify(notification *dto.Notification)
	SetHandler(method string, handler WsHandler)
	AddClient(ticketer Ticketer, dal DAL, w http.ResponseWriter, r *http.Request)
	Run()
}
