package services

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/hasher"
)

type Authenticator struct {
	users []config.User
}

func (authenticator *Authenticator) getUser(username string) *config.User {
	for _, user := range authenticator.users {
		if user.Username == username {
			return &user
		}
	}

	return nil
}

func (authenticator *Authenticator) AuthenticateFromToken(username string, salt string, token string) bool {
	user := authenticator.getUser(username)
	if user == nil {
		return false
	}

	return token == hasher.GetHash([]byte(user.Password+salt))
}

func (authenticator *Authenticator) AuthenticateFromPassword(username string, password string) bool {
	user := authenticator.getUser(username)
	if user == nil {
		return false
	}

	return user.Password == password
}

func (authenticator *Authenticator) AuthenticateFromJWT(token string) {

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		// Check the hashing algorithm is HMAC
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// Provide the hashing secret so the token claims can be verified
		return []byte("my_secret_key"), nil
	})

	claims, ok := t.Claims.(jwt.MapClaims)
	if ok && t.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}
}

func NewAuthenticator(config *config.Config) *Authenticator {
	return &Authenticator{
		users: config.Users,
	}
}
