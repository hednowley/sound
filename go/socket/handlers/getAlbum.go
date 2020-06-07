package handlers

import (
	"encoding/json"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetAlbumHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request, _ *socket.HandlerContext) interface{} {
		var id uint

		if request.Params["id"] == nil || json.Unmarshal(*request.Params["id"], &id) != nil {
			return "bad id"
		}

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse("Cannot make DB conn")
		}
		defer conn.Release()

		album, err := dal.Db.GetAlbum(conn, id)
		if err != nil {
			return "no album"
		}

		songs, err := dal.Db.GetAlbumSongs(conn, id)
		if err != nil {
			return api.NewErrorReponse("Error")
		}

		return dto.NewAlbum(album, songs)
	}
}
