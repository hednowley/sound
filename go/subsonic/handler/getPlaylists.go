package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewGetPlaylistsHandler does http://www.subsonic.org/pages/api.jsp#getPlaylists
func NewGetPlaylistsHandler(dal interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		playlists := dal.GetPlaylists()
		return api.NewSuccessfulReponse(dto.NewPlaylistCollection(playlists))
	}
}
