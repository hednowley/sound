package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewCreatePlaylistHandler does http://www.subsonic.org/pages/api.jsp#createPlaylist
func NewCreatePlaylistHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values, context *api.HandlerContext) *api.Response {

		idParam := params.Get("playlistId")
		id := util.ParseUint(idParam, 0)
		name := params.Get("name")

		if id == 0 && name == "" {
			return api.NewErrorReponse(dto.MissingParameter, "Playlist ID or name must be specified.")
		}

		songIds := params["songId"]

		songIdsNum := []uint{}
		for _, idStr := range songIds {
			id := util.ParseUint(idStr, 0)
			if id != 0 {
				songIdsNum = append(songIdsNum, id)
			}
		}

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		if id == 0 {
			// We are being asked to create a playlist
			inserted, err := dal.Db.InsertPlaylist(conn, name, "", context.User.Username, false)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}

			// Use the created playlist's ID
			id = inserted
		} else {
			// We are being asked to update a playlist
			playlist, err := dal.Db.GetPlaylist(conn, id, context.User.Username)
			if err != nil {
				_, ok := err.(*dao.ErrNotFound)
				if ok {
					m := fmt.Sprintf("Playlist not found: %v", idParam)
					return api.NewErrorReponse(dto.NotFound, m)
				}
				return api.NewErrorReponse(dto.Generic, err.Error())
			}

			var nameUpdate string
			if name == "" {
				nameUpdate = playlist.Name
			} else {
				nameUpdate = name
			}

			playlist, err = dal.Db.UpdatePlaylist(conn, playlist.ID, nameUpdate, playlist.Public, playlist.Comment)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}
		}

		// Set the playlist songs
		err = dal.Db.ReplacePlaylistEntries(conn, id, songIdsNum)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		// Fetch the updated playlist
		p, err := dal.Db.GetPlaylist(conn, id, context.User.Username)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		songs, err := dal.Db.GetPlaylistSongs(conn, id, context.User.Username)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewPlaylist(p, songs))
	}
}
