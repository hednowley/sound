package handler

import (
	"net/url"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetUserHandler(config *config.Config) api.Handler {

	return func(params url.Values) *api.Response {

		username := params.Get("username")
		if len(username) == 0 {
			return api.NewErrorReponse(dto.NotFound, "No username.")
		}

		for _, user := range config.Users {
			if user.Username == username {
				return api.NewSuccessfulReponse(dto.NewUser(user))
			}
		}

		return api.NewErrorReponse(dto.NotFound, "User not found.")
	}
}
