package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

// NewDeletePlaylistHandler is a handler for deleting playlists.
func NewDeletePlaylistHandler(database *dao.Database) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("id")
		id := api.ParseUint(idParam, 0)
		if id == 0 {
			message := fmt.Sprintf("Playlist not found: %v", idParam)
			return api.NewErrorReponse(dto.NotFound, message)
		}

		database.DeletePlaylist(id)
		return api.NewSuccessfulReponse(nil)
	}
}
