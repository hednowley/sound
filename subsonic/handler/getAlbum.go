package handler

import (
	"net/url"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewGetAlbumHandler is a handler for getting information about an album.
func NewGetAlbumHandler(database interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.MissingParameter, "Required param (id) is missing")
		}

		album, err := database.GetAlbum(id, false, false, true)
		if err != nil {
			if _, ok := err.(*dao.ErrNotFound); ok {
				return api.NewErrorReponse(dto.NotFound, "Album not found.")
			}
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewAlbum(album, true))
	}
}
