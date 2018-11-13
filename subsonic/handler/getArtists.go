package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetArtistsHandler(database *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		artists := database.GetArtists()
		return api.NewSuccessfulReponse(dto.NewArtistCollection(artists))
	}
}
