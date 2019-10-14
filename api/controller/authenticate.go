package controller

import (
	"net/http"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

// NewAuthenticateController makes a controller which gives out JWT tokens in return for credentials.
func NewAuthenticateController(authenticator *services.Authenticator) *api.BinaryController {

	input := dto.Credentials{}

	run := func(w *http.ResponseWriter, r *http.Request, _ *config.User) *api.Response {
		credentials := &input
		if !authenticator.AuthenticateFromPassword(credentials.Username, credentials.Password) {
			return api.NewErrorReponse("Bad credentials.")
		}

		token, err := authenticator.MakeJWT(credentials.Username)
		if err != nil {
			return api.NewErrorReponse("Could not make token.")
		}

		cookie := http.Cookie{
			Name:  "access_token",
			Value: token,
		}
		http.SetCookie(*w, &cookie)

		return api.NewSuccessfulReponse(&struct{}{})
	}

	return &api.BinaryController{
		Input:  &input,
		Run:    run,
		Secure: false,
	}
}
