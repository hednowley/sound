package handlers

import "github.com/hednowley/sound/ws/dto"

// WsHandler listens for particular websocket messages.
type WsHandler = func(*dto.Request) interface{}
