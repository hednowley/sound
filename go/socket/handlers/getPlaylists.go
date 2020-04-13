package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetPlaylistsHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request) interface{} {
		playlists := dal.Db.GetPlaylists()
		return dto.NewPlaylistCollection(playlists)
	}
}
