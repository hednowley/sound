package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewGetGenresHandler does http://www.subsonic.org/pages/api.jsp#getGenres
func NewGetGenresHandler(dal interfaces.DAL) api.Handler {
	return func(params url.Values) *api.Response {
		genres := dal.GetGenres()
		return api.NewSuccessfulReponse(dto.NewGenres(genres))
	}
}
