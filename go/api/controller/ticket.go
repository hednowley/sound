package controller

import (
	"net/http"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/socket"
)

// NewTicketController makes a controller which returns a new Websocket ticket.
func NewTicketController(ticketer *socket.Ticketer) *api.Controller {

	make := func() *api.ControllerContext {

		return &api.ControllerContext{
			Body: nil,
			Run: func(user *config.User, _ http.ResponseWriter, _ *http.Request) *api.Response {
				r := dto.NewTicket(ticketer.MakeTicket(user))
				return api.NewSuccessfulReponse(&r)
			},
		}
	}

	return &api.Controller{
		Secure: true,
		Make:   make,
	}
}
