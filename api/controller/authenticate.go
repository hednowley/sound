package controller

import (
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

// NewAuthenticateController makes a controller which gives out JWT tokens.
func NewAuthenticateController(authenticator *services.Authenticator) *api.Controller {

	input := dto.Credentials{}

	w := func(_ *config.User) *api.Response {
		credentials := &input
		if !authenticator.AuthenticateFromPassword(credentials.Username, credentials.Password) {
			return api.NewErrorReponse("Bad credentials.")
		}

		token, err := authenticator.MakeJWT(credentials.Username)
		if err != nil {
			return api.NewErrorReponse(err.Error())
		}

		return api.NewSuccessfulReponse(&dto.Token{Token: token})
	}

	return &api.Controller{
		Input:  &input,
		Run:    w,
		Secure: false,
	}
}
