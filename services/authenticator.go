package services

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/hasher"
)

type Authenticator struct {
	users  []config.User
	secret string
}

func NewAuthenticator(config *config.Config) *Authenticator {
	return &Authenticator{
		users:  config.Users,
		secret: config.Secret,
	}
}

func (a *Authenticator) getUser(username string) *config.User {
	for _, user := range a.users {
		if user.Username == username {
			return &user
		}
	}

	return nil
}

// These methods should probably return a claims struct rather than bool...

func (a *Authenticator) AuthenticateFromToken(username string, salt string, token string) bool {
	user := a.getUser(username)
	if user == nil {
		return false
	}

	return token == hasher.GetHash([]byte(user.Password+salt))
}

func (a *Authenticator) AuthenticateFromPassword(username string, password string) bool {
	user := a.getUser(username)
	if user == nil {
		return false
	}

	return user.Password == password
}

func (a *Authenticator) AuthenticateFromJWT(token string) bool {

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
		return false
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if ok && t.Valid {
		u, ok := claims["u"].(string)
		if !ok {
			return false
		}
		return a.getUser(u) != nil
	} else {
		return false
	}
}
