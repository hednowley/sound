package controller

import (
	"net/http"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

// NewAuthenticateController makes a controller which gives out JWT tokens in return for credentials.
func NewAuthenticateController(authenticator *services.Authenticator) *api.Controller {

	make := func() *api.ControllerContext {

		credentials := &dto.Credentials{}

		run := func(_ *config.User, w http.ResponseWriter, r *http.Request) *api.Response {

			if authenticator.AuthenticateFromPassword(credentials.Username, credentials.Password) == nil {
				return api.NewErrorReponse("Bad credentials.")
			}

			token, err := authenticator.MakeJWT(credentials.Username)
			if err != nil {
				return api.NewErrorReponse("Could not make token.")
			}

			//expire := time.Now().AddDate(0, 0, 1)
			cookie := http.Cookie{
				Name:  "token",
				Value: token,
				//Domain:  "false",
				//Expires: expire,
			}
			http.SetCookie(w, &cookie)

			return api.NewSuccessfulReponse(&struct{}{})
		}

		return &api.ControllerContext{
			Body: credentials,
			Run:  run,
		}
	}

	return &api.Controller{
		Make:   make,
		Secure: false,
	}
}
