package handlers

import (
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetPlaylistsHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request) interface{} {
		playlists, err := dal.Db.GetPlaylists()
		if err != nil {
			return api.NewErrorReponse("Error")
		}

		return dto.NewPlaylistCollection(playlists)
	}
}
