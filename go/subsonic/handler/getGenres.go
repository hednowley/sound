package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewGetGenresHandler does http://www.subsonic.org/pages/api.jsp#getGenres
func NewGetGenresHandler(dal *dal.DAL) api.Handler {
	return func(params url.Values) *api.Response {
		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		genres, err := dal.Db.GetGenres(conn)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewGenres(genres))
	}
}
