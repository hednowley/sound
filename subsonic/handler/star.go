package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewStarHandler(dal interfaces.DAL, star bool) api.Handler {

	return func(params url.Values) *api.Response {

		param := params.Get("id")
		id := api.ParseUint(param, 0)
		if id == 0 {
			err := dal.StarSong(id, star)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}

			return api.NewEmptyReponse()
		}

		param = params.Get("albumId")
		id = api.ParseUint(param, 0)
		if id == 0 {
			err := dal.StarAlbum(id, star)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}

			return api.NewEmptyReponse()
		}

		param = params.Get("artistId")
		id = api.ParseUint(param, 0)
		if id == 0 {
			err := dal.StarArtist(id, star)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}

			return api.NewEmptyReponse()
		}

		return api.NewErrorReponse(dto.Generic, "Missing parameter")
	}
}