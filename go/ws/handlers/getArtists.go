package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/ws"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetArtistsHandler(dal *dal.DAL) ws.WsHandler {
	return func(request *dto.Request) interface{} {
		artists := dal.GetArtists(false)
		return dto.NewArtistCollection(artists)
	}
}
