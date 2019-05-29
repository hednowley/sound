package handlers

import (
	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/ws"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetArtistsHandler(dal idal.DAL) ws.WsHandler {
	return func(request *dto.Request) interface{} {
		artists := dal.GetArtists()
		return dto.NewArtistCollection(artists)
	}
}
