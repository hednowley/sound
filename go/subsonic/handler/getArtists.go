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
		artists, err := dal.Db.GetArtists()
		if err != nil {
			return api.NewErrorReponse(0, "Error")
		}
		return api.NewSuccessfulReponse(dto.NewArtistCollection(artists, conf))
	}
}
