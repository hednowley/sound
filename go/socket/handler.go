package socket

import (
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/socket/dto"
)

type HandlerContext struct {
	User *config.User
}

// socket.Handler listens for particular websocket messages and
// returns an object which will be sent back to the sender.
type Handler = func(*dto.Request, *HandlerContext) interface{}
