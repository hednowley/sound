package handler

import (
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

func NewGetMusicDirectoryHandler(database *dao.Database) api.Handler {

	return func(params url.Values) *api.Response {

		id := api.ParseUint(params.Get("id"), 0)
		if id == 0 {
			return api.NewErrorReponse(dto.Generic, "Bad id")
		}

		album, err := database.GetAlbum(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewDirectory(album))
	}
}
