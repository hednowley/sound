package handler

import (
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dto"
)

func NewGetUsersHandler(config *config.Config) api.Handler {

	return func(params url.Values) *api.Response {
		return api.NewSuccessfulReponse(dto.NewUserCollection(config.Users))
	}
}
