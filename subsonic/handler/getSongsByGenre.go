package handler

import (
	"net/url"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewGetSongsByGenreHandler(dal interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		genreParam := params.Get("genre")
		if len(genreParam) == 0 {
			return api.NewErrorReponse(dto.MissingParameter, "Required param (genre) is missing")
		}

		countParam := params.Get("count")
		count := util.ParseUint(countParam, 10)

		offsetParam := params.Get("offset")
		offset := util.ParseUint(offsetParam, 0)

		genre, err := dal.GetGenre(genreParam)
		if err != nil {
			genre = &dao.Genre{}
		}

		return api.NewSuccessfulReponse(dto.NewSongsByGenre(genre, count, offset))
	}
}
