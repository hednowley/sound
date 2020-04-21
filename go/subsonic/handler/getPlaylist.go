package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewGetPlaylistHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			message := fmt.Sprintf("Playlist not found: %v", idParam)
			return api.NewErrorReponse(dto.Generic, message)
		}

		p, err := dal.Db.GetPlaylist(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		songs, err := dal.Db.GetPlaylistSongs(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewPlaylist(p, songs))
	}
}
