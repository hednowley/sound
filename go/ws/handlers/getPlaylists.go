package handlers

import (
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetPlaylistsHandler(dal interfaces.DAL) interfaces.WsHandler {
	return func(request *dto.Request) interface{} {
		playlists := dal.GetPlaylists()
		return dto.NewPlaylistCollection(playlists)
	}
}
