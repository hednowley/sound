package controller

import (
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/ws"
)

// NewTicketController makes a controller which returns a new Websocket ticket.
func NewTicketController(ticketer *ws.Ticketer) *api.Controller {

	input := struct{}{}

	w := func(user *config.User) *api.Response {
		r := dto.NewTicket(ticketer.MakeTicket(user))
		return api.NewSuccessfulReponse(&r)
	}

	return &api.Controller{
		Input:  &input,
		Run:    w,
		Secure: true,
	}
}
