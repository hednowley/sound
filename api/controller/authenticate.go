package controller

import (
	"encoding/json"
	"net/http"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

// NewAuthenticateController makes a controller which gives out JWT tokens in return for credentials.
func NewAuthenticateController(authenticator *services.Authenticator) *api.BinaryController {

	run := func(w *http.ResponseWriter, r *http.Request, _ *config.User) *api.Response {

		credentials := &dto.Credentials{}

		if r.Body == http.NoBody {
			return api.NewErrorReponse("Bad request.")
		}

		d := json.NewDecoder(r.Body)
		err := d.Decode(credentials)
		if err != nil {
			return api.NewErrorReponse("Bad request.")
		}

		if !authenticator.AuthenticateFromPassword(credentials.Username, credentials.Password) {
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
		http.SetCookie(*w, &cookie)

		return api.NewSuccessfulReponse(&struct{}{})
	}

	return &api.BinaryController{
		Run:    run,
		Secure: false,
	}
}
