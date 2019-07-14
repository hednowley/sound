package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetIndexesHandler(database interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		artists := database.GetArtists()
		return api.NewSuccessfulReponse(dto.NewIndexes(artists))
	}
}
