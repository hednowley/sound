package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetArtistsHandler(dal *dal.DAL) WsHandler {
	return func(request *dto.Request) interface{} {
		artists := dal.GetArtists()
		return dto.NewArtistCollection(artists)
	}
}
