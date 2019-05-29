package handler

import (
	"net/url"

	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetGenresHandler(database idal.DAL) api.Handler {
	return func(params url.Values) *api.Response {
		genres := database.GetGenres()
		return api.NewSuccessfulReponse(dto.NewGenres(genres))
	}
}
