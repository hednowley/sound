package interfaces

import "github.com/hednowley/sound/ws/dto"

type Hub interface {
	Notify(notification *dto.Notification)
}
