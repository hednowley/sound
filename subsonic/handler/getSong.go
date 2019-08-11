package handler

import (
	"net/url"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewGetSongHandler(database interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.MissingParameter, "Required param (id) is missing")
		}

		file, err := database.GetSong(id, true, true, true)
		if err != nil {
			if dao.IsErrNotFound(err) {
				return api.NewErrorReponse(dto.NotFound, "Song not found.")
			}
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewSong(file))
	}
}
