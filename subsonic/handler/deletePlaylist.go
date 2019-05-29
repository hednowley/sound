package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewDeletePlaylistHandler is a handler for deleting playlists.
func NewDeletePlaylistHandler(database idal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("id")
		id := api.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.MissingParameter, "Required param (id) is missing")
		}

		err := database.DeletePlaylist(id)
		if err != nil {
			if _, ok := err.(*dao.ErrNotFound); ok {
				message := fmt.Sprintf("Playlist not found: %v", idParam)
				return api.NewErrorReponse(dto.NotFound, message)
			}
			return api.NewErrorReponse(dto.NotFound, err.Error())
		}

		return api.NewSuccessfulReponse(nil)
	}
}
