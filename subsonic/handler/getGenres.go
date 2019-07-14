package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetGenresHandler(database interfaces.DAL) api.Handler {
	return func(params url.Values) *api.Response {
		genres := database.GetGenres()
		return api.NewSuccessfulReponse(dto.NewGenres(genres))
	}
}
