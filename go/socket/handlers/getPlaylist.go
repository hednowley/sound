package handlers

import (
	"encoding/json"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetPlaylistHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request, context *socket.HandlerContext) interface{} {
		var id uint

		if request.Params["id"] == nil || json.Unmarshal(*request.Params["id"], &id) != nil {
			return "bad id"
		}

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse("Cannot make DB conn")
		}
		defer conn.Release()

		playlist, err := dal.Db.GetPlaylist(conn, id, context.User.Username)
		if err != nil {
			return "no playlist"
		}

		songs, err := dal.Db.GetPlaylistSongs(conn, id, context.User.Username)
		if err != nil {
			return "no playlist"
		}

		return dto.NewPlaylist(playlist, songs)
	}
}
