package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/config"
)

type UserCollection struct {
	XMLName xml.Name `xml:"users" json:"-"`
	Users   []User   `xml:"user" json:"user"`
}

func NewUserCollection(users []config.User) UserCollection {

	col := make([]User, len(users))
	for index, user := range users {
		col[index] = NewUser(user)
	}

	return UserCollection{
		Users: col,
	}
}
