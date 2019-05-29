package handler

import (
	"net/url"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetSongsByGenreHandler(dal idal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		genreParam := params.Get("genre")
		if len(genreParam) == 0 {
			return api.NewErrorReponse(dto.MissingParameter, "Required param (genre) is missing")
		}

		countParam := params.Get("count")
		count := api.ParseUint(countParam, 10)

		offsetParam := params.Get("offset")
		offset := api.ParseUint(offsetParam, 0)

		genre, err := dal.GetGenre(genreParam)
		if err != nil {
			genre = &dao.Genre{}
		}

		return api.NewSuccessfulReponse(dto.NewSongsByGenre(genre, count, offset))
	}
}
