package services

import (
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

func NewAuthenticator(config *config.Config) *Authenticator {
	return &Authenticator{
		users: config.Users,
	}
}
