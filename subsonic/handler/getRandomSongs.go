package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewGetRandomSongsHandler(database interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		sizeParam := params.Get("size")
		size := util.ParseUint(sizeParam, 10)

		genre := params.Get("genre")

		fromParam := params.Get("fromYear")
		from := util.ParseUint(fromParam, 0)

		toParam := params.Get("toYear")
		to := util.ParseUint(toParam, 0)

		songs := database.GetRandomSongs(size, from, to, genre)
		return api.NewSuccessfulReponse(dto.NewRandomSongs(songs))
	}
}
