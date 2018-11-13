package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/subsonic/api"
)

func NewGetGenresHandler(database *dal.DAL) api.Handler {
	return func(params url.Values) *api.Response {
		genres := database.GetGenres()
		return api.NewSuccessfulReponse(dto.NewGenres(genres))
	}
}
