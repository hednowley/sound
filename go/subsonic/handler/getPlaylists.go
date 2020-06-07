package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewGetPlaylistsHandler does http://www.subsonic.org/pages/api.jsp#getPlaylists
func NewGetPlaylistsHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values, context *api.HandlerContext) *api.Response {
		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		playlists, err := dal.Db.GetPlaylists(conn, context.User.Username)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		// Pretend that all accessible playlists are owned by the requestor.
		// Allows anyone to edit a public playlist (against the intention of
		// the Subsonic API).
		for _, p := range playlists {
			p.Owner = context.User.Username
		}

		return api.NewSuccessfulReponse(dto.NewPlaylistCollection(playlists))
	}
}
