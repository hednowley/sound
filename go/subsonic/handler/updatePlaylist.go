package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewUpdatePlaylistHandler does http://www.subsonic.org/pages/api.jsp#updatePlaylist
func NewUpdatePlaylistHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values, context *api.HandlerContext) *api.Response {

		idParam := params.Get("playlistId")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			message := fmt.Sprintf("Playlist not found: %v", idParam)
			return api.NewErrorReponse(dto.NotFound, message)
		}

		name := params.Get("name")
		comment := params.Get("comment")

		public := util.ParseBool(params.Get("public"))

		addedSongsParam := params["songIdToAdd"]
		addedSongs := []uint{}
		for _, idStr := range addedSongsParam {
			songID := util.ParseUint(idStr, 0)
			if songID != 0 {
				addedSongs = append(addedSongs, songID)
			}
		}

		removedSongsParam := params["songIndexToRemove"]
		removedSongs := []uint{}
		for _, idStr := range removedSongsParam {
			songID := util.ParseUint(idStr, 0)
			if songID != 0 {
				removedSongs = append(removedSongs, songID)
			}
		}

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		err = dal.UpdatePlaylist(conn, id, name, comment, public, addedSongs, removedSongs, context.User.Username)
		if err != nil {
			return api.NewErrorReponse(0, err.Error())
		}

		return api.NewEmptyReponse()
	}
}
