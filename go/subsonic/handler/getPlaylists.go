package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewGetPlaylistsHandler does http://www.subsonic.org/pages/api.jsp#getPlaylists
func NewGetPlaylistsHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		playlists, err := dal.Db.GetPlaylists()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewPlaylistCollection(playlists))
	}
}
