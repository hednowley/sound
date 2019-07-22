package controller

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

// NewAuthenticateController makes a controller which gives out JWT tokens.
func NewAuthenticateController(authenticator *services.Authenticator, cfg *config.Config) *api.Controller {

	input := dto.Credentials{}

	w := func(_ *config.User) *api.Response {
		credentials := &input
		if !authenticator.AuthenticateFromPassword(credentials.Username, credentials.Password) {
			return api.NewErrorReponse("Bad credentials.")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"u": credentials.Username,
			//"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, _ := token.SignedString([]byte(cfg.Secret))
		return api.NewSuccessfulReponse(&dto.Token{Token: tokenString})
	}

	return &api.Controller{
		Input:  &input,
		Run:    w,
		Secure: false,
	}
}
