package handler

import (
	"net/url"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetUsersHandler(config *config.Config) api.Handler {

	return func(params url.Values, _ *api.HandlerContext) *api.Response {
		return api.NewSuccessfulReponse(dto.NewUserCollection(config.Users))
	}
}
