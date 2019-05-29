package handler

import (
	"net/url"

	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetArtistsHandler(database idal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		artists := database.GetArtists()
		return api.NewSuccessfulReponse(dto.NewArtistCollection(artists))
	}
}
