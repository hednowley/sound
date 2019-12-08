package ws

import (
	"github.com/hednowley/sound/ws/dto"
)

// ws.WsHandler listens for particular websocket messages and
// returns an object which will be sent back to the sender.
type WsHandler = func(*dto.Request) interface{}
