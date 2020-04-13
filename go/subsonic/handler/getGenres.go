package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewGetGenresHandler does http://www.subsonic.org/pages/api.jsp#getGenres
func NewGetGenresHandler(dal *dal.DAL) api.Handler {
	return func(params url.Values) *api.Response {
		genres := dal.Db.GetGenres()
		return api.NewSuccessfulReponse(dto.NewGenres(genres))
	}
}
