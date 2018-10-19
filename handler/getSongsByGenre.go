package handler

import (
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

func NewGetSongsByGenreHandler(database *dao.Database) api.Handler {

	return func(params url.Values) *api.Response {

		genreParam := params.Get("genre")

		countParam := params.Get("count")
		count := api.ParseUint(countParam, 10)

		offsetParam := params.Get("offset")
		offset := api.ParseUint(offsetParam, 0)

		genre, err := database.GetGenre(genreParam)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewSongsByGenre(genre, count, offset))
	}
}
