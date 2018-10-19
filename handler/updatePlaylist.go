package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

func NewUpdatePlaylistHandler(database *dao.Database) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("playlistId")
		id := api.ParseUint(idParam, 0)
		if id == 0 {
			message := fmt.Sprintf("Playlist not found: %v", idParam)
			return api.NewErrorReponse(dto.NotFound, message)
		}

		name := params.Get("name")
		comment := params.Get("comment")

		public := api.ParseBool(params.Get("public"))

		addedSongsParam := params["songIdToAdd"]
		addedSongs := []uint{}
		for _, idStr := range addedSongsParam {
			songID := api.ParseUint(idStr, 0)
			if songID != 0 {
				addedSongs = append(addedSongs, songID)
			}
		}

		removedSongsParam := params["songIndexToRemove"]
		removedSongs := []uint{}
		for _, idStr := range removedSongsParam {
			songID := api.ParseUint(idStr, 0)
			if songID != 0 {
				removedSongs = append(removedSongs, songID)
			}
		}

		err := database.UpdatePlaylist(id, name, comment, public, addedSongs, removedSongs)
		if err != nil {
			api.NewErrorReponse(0, err.Error())
		}

		return &api.Response{
			Body:      nil,
			IsSuccess: true,
		}
	}
}
