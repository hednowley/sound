package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewStarHandler(dal *dal.DAL, star bool) api.Handler {

	return func(params url.Values, _ *api.HandlerContext) *api.Response {

		param := params.Get("id")
		id := util.ParseUint(param, 0)
		if id == 0 {
			err := dal.StarSong(id, star)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}

			return api.NewEmptyReponse()
		}

		param = params.Get("albumId")
		id = util.ParseUint(param, 0)
		if id == 0 {
			err := dal.StarAlbum(id, star)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}

			return api.NewEmptyReponse()
		}

		param = params.Get("artistId")
		id = util.ParseUint(param, 0)
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
