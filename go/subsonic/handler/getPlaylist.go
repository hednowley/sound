package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewGetPlaylistHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values, context *api.HandlerContext) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			message := fmt.Sprintf("Playlist not found: %v", idParam)
			return api.NewErrorReponse(dto.Generic, message)
		}

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		p, err := dal.Db.GetPlaylist(conn, id, context.User.Username)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		// Pretend that this playlist is owned by the requestor.
		// Allows anyone to edit a public playlist (against the intention of
		// the Subsonic API).
		p.Owner = context.User.Username

		songs, err := dal.Db.GetPlaylistSongs(conn, id, context.User.Username)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewPlaylist(p, songs))
	}
}
