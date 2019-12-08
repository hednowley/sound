package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/ws"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetPlaylistsHandler(dal *dal.DAL) ws.WsHandler {
	return func(request *dto.Request) interface{} {
		playlists := dal.GetPlaylists()
		return dto.NewPlaylistCollection(playlists)
	}
}
