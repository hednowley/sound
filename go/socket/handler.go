package socket

import (
	"github.com/hednowley/sound/socket/dto"
)

// socket.Handler listens for particular websocket messages and
// returns an object which will be sent back to the sender.
type Handler = func(*dto.Request) interface{}
