package handlers

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

// NewAuthenticateHandler makes a handler which accepts credentials and returns
// a JWT token if they are valid.
func NewAuthenticateHandler(config *config.Config) *api.Controller {

	input := dto.Credentials{}

	w := func() *api.Response {
		credentials := &input
		a := services.NewAuthenticator(config)
		if !a.AuthenticateFromPassword(credentials.Username, credentials.Password) {
			return api.NewErrorReponse("Bad credentials.")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"u": credentials.Username,
			//"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, _ := token.SignedString([]byte(config.Secret))
		return api.NewSuccessfulReponse(dto.Token{Token: tokenString})
	}

	return &api.Controller{
		Input:  &input,
		Run:    w,
		Secure: false,
	}
}
