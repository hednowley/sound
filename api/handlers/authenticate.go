package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

// NewAuthenticateHandler makes a handler which accepts credentials and returns
// a JWT token if they are valid.
func NewAuthenticateHandler(config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c := new(dto.Credentials)
		d := json.NewDecoder(r.Body)
		err := d.Decode(&c)
		if err != nil {
			s := dto.MakeResponse(dto.Fail, "Malformed credentials.")
			fmt.Fprint(w, s)
			return
		}

		a := services.NewAuthenticator(config)
		if !a.AuthenticateFromPassword(c.Username, c.Password) {
			s := dto.MakeResponse(dto.Fail, "Bad credentials.")
			fmt.Fprint(w, s)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"u": c.Username,
			//"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(config.Secret))
		s := dto.MakeResponse(dto.Success, dto.Token{Token: tokenString})
		fmt.Fprint(w, s)
	}
}
