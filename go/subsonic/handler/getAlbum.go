package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewGetAlbumHandler does http://www.subsonic.org/pages/api.jsp#getAlbum
func NewGetAlbumHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

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

		album, err := dal.Db.GetAlbum(conn, id)
		if err != nil {
			if _, ok := err.(*dao.ErrNotFound); ok {
				return api.NewErrorReponse(dto.NotFound, "Album not found.")
			}
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		songs, err := dal.Db.GetAlbumSongs(conn, id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewAlbumWithSongs(album, songs))
	}
}
