package handler

import (
	"net/url"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetArtistsHandler(dal *dal.DAL, conf *config.Config) api.Handler {
	return func(params url.Values) *api.Response {
		artists := dal.Db.GetArtists()
		return api.NewSuccessfulReponse(dto.NewArtistCollection(artists, conf))
	}
}
