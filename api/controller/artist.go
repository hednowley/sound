package controller

import (
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
)

func NewArtistController() *api.Controller {

	input := struct{}{}

	w := func() *api.Response {
		return api.NewSuccessfulReponse(&dto.Token{Token: "hello"})
	}

	return &api.Controller{
		Input:  &input,
		Run:    w,
		Secure: true,
	}
}
