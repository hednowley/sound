package controller

import (
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/dal"
)

func NewArtistCollectionController(dal *dal.DAL) *api.Controller {

	input := struct{}{}

	w := func() *api.Response {
		artists := dal.GetArtists()
		return api.NewSuccessfulReponse(dto.NewArtistCollection(artists))
	}

	return &api.Controller{
		Input:  &input,
		Run:    w,
		Secure: true,
	}
}
