package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewGetArtistHandler does http://www.subsonic.org/pages/api.jsp#getArtist
func NewGetArtistHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values, _ *api.HandlerContext) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.MissingParameter, "Required param (id) is missing")
		}

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		artist, err := dal.Db.GetArtist(conn, id)
		if err != nil {
			if _, ok := err.(*dao.ErrNotFound); ok {
				return api.NewErrorReponse(dto.NotFound, "Artist not found.")
			}
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		albums, err := dal.Db.GetAlbumsByArtist(conn, id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewArtistWithAlbums(artist, albums))
	}
}
