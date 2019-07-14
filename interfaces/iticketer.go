package interfaces

import "github.com/hednowley/sound/config"

type Ticketer interface {
	SubmitTicket(key string) *config.User
}
