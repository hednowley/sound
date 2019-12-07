package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetArtistsHandler(dal *dal.DAL) interfaces.WsHandler {
	return func(request *dto.Request) interface{} {
		artists := dal.GetArtists(false)
		return dto.NewArtistCollection(artists)
	}
}
