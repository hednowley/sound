package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewGetSongsByGenreHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values, _ *api.HandlerContext) *api.Response {

		genre := params.Get("genre")
		if len(genre) == 0 {
			return api.NewErrorReponse(dto.MissingParameter, "Required param (genre) is missing")
		}

		countParam := params.Get("count")
		count := util.ParseUint(countParam, 10)

		offsetParam := params.Get("offset")
		offset := util.ParseUint(offsetParam, 0)

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		songs, err := dal.Db.GetSongsByGenre(conn, genre, offset, count)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewSongsByGenre(songs))
	}
}
