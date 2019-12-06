package handler

import (
	"net/url"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetArtistsHandler(dal interfaces.DAL, conf *config.Config) api.Handler {
	return func(params url.Values) *api.Response {
		artists := dal.GetArtists(false)
		return api.NewSuccessfulReponse(dto.NewArtistCollection(artists, conf))
	}
}
