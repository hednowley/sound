package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetPlaylistsHandler(database interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		playlists := database.GetPlaylists()
		return api.NewSuccessfulReponse(dto.NewPlaylistCollection(playlists))
	}
}
