package handlers

import "github.com/hednowley/sound/ws/dto"

type WsHandler = func(*dto.Request) interface{}
