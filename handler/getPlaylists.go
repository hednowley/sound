package handler

import (
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

func NewGetPlaylistsHandler(database *dao.Database) api.Handler {

	return func(params url.Values) *api.Response {

		playlists := database.GetPlaylists()
		return api.NewSuccessfulReponse(dto.NewPlaylistCollection(playlists))
	}
}
