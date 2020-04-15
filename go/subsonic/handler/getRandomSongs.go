package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewGetRandomSongsHandler does http://www.subsonic.org/pages/api.jsp#getRandomSongs
func NewGetRandomSongsHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		sizeParam := params.Get("size")
		size := util.ParseUint(sizeParam, 10)

		genre := params.Get("genre")

		fromParam := params.Get("fromYear")
		from := util.ParseUint(fromParam, 0)

		toParam := params.Get("toYear")
		to := util.ParseUint(toParam, 0)

		songs, err := dal.Db.GetRandomSongs(size, from, to, genre)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewRandomSongs(songs))
	}
}
