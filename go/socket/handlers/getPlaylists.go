package handlers

import (
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetPlaylistsHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request, context *socket.HandlerContext) interface{} {
		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse("Cannot make DB conn")
		}
		defer conn.Release()

		playlists, err := dal.Db.GetPlaylists(conn, context.User.Username)
		if err != nil {
			return api.NewErrorReponse("Error")
		}

		return dto.NewPlaylistCollection(playlists)
	}
}
