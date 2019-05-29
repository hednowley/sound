package handler

import (
	"net/url"

	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetPlaylistsHandler(database idal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		playlists := database.GetPlaylists()
		return api.NewSuccessfulReponse(dto.NewPlaylistCollection(playlists))
	}
}
