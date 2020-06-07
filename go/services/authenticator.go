package services

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/hasher"
)

// Authenticator verifies authentication claims.
// TODO: Return claims structs rather than bools
type Authenticator struct {
	users  []config.User
	secret string
}

// NewAuthenticator constructs a new authenticator against the configured user-set.
func NewAuthenticator(config *config.Config) *Authenticator {
	return &Authenticator{
		users:  config.Users,
		secret: config.Secret,
	}
}

// getUser tries to retrieve the user with the given name.
func (a *Authenticator) getUser(username string) *config.User {
	for _, user := range a.users {
		if user.Username == username {
			return &user
		}
	}

	return nil
}

// AuthenticateFromToken verifies credentials where the password has been salted and hashed
// in the format expected by Subsonic.
func (a *Authenticator) AuthenticateFromToken(username string, salt string, token string) *config.User {
	user := a.getUser(username)
	if user == nil || token != hasher.GetHash([]byte(user.Password+salt)) {
		return nil
	}

	return user
}

// AuthenticateFromPassword verifies credentials where the password has been provided in plain text
// per the deprecated Subsonic API.
func (a *Authenticator) AuthenticateFromPassword(username string, password string) *config.User {
	user := a.getUser(username)
	if user == nil || user.Password != password {
		return nil
	}

	return user
}

// AuthenticateFromJWT verifies auth claims encoded as a JSON web token.
func (a *Authenticator) AuthenticateFromJWT(token string) *config.User {

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		// Check the hashing algorithm is HMAC
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// Provide the hashing secret so the token claims can be verified
		return []byte(a.secret), nil
	})

	if err != nil {
		return nil
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if ok && t.Valid {
		u, ok := claims["u"].(string)
		if !ok {
			return nil
		}
		return a.getUser(u)
	}
	return nil
}

// MakeJWT creates a JSON web token from an auth claim.
func (a *Authenticator) MakeJWT(username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"u": username,
		//"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	return token.SignedString([]byte(a.secret))
}
