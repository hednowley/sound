package handler

import (
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

func NewGetIndexesHandler(database *dao.Database) api.Handler {

	return func(params url.Values) *api.Response {
		artists := database.GetArtists()
		return api.NewSuccessfulReponse(dto.NewIndexes(artists))
	}
}
