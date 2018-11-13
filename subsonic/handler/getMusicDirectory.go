package handler

import (
	"net/url"

	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetMusicDirectoryHandler(database *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		id := api.ParseUint(params.Get("id"), 0)
		if id == 0 {
			return api.NewErrorReponse(dto.Generic, "Bad id")
		}

		album, err := database.GetAlbum(id, true, true, true)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewDirectory(album))
	}
}
