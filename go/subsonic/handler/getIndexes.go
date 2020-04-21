package handler

import (
	"net/url"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewGetIndexesHandler does http://www.subsonic.org/pages/api.jsp#getIndexes
func NewGetIndexesHandler(dal *dal.DAL, conf *config.Config) api.Handler {
	return func(params url.Values) *api.Response {
		artists, err := dal.Db.GetArtists()
		if err != nil {
			return api.NewErrorReponse(0, "Error")
		}
		return api.NewSuccessfulReponse(dto.NewIndexCollection(artists, conf))
	}
}
