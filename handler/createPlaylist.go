package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

// NewCreatePlaylistHandler is a handler for creating or updating playlists.
func NewCreatePlaylistHandler(database *dao.Database) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("playlistId")
		id := api.ParseUint(idParam, 0)
		name := params.Get("name")

		if id == 0 && name == "" {
			return api.NewErrorReponse(dto.MissingParameter, "Playlist ID or name must be specified.")
		}

		songIds := params["songId"]

		songIdsNum := []uint{}
		for _, idStr := range songIds {
			id := api.ParseUint(idStr, 0)
			if id != 0 {
				songIdsNum = append(songIdsNum, id)
			}
		}

		id, err := database.PutPlaylist(id, name, songIdsNum)
		if err != nil {
			_, ok := err.(*dao.ErrNotFound)
			if ok {
				m := fmt.Sprintf("Playlist not found: %v", idParam)
				return api.NewErrorReponse(dto.NotFound, m)
			}
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		playlist, err := database.GetPlaylist(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewPlaylist(playlist))
	}
}
