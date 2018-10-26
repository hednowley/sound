package handler

import (
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dto"
)

func NewGetGenresHandler(database *dal.DAL) api.Handler {
	return func(params url.Values) *api.Response {
		genres := database.GetGenres()
		return api.NewSuccessfulReponse(dto.NewGenres(genres))
	}
}
