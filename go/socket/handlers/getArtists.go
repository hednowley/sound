package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetArtistsHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request) interface{} {
		artists := dal.Db.GetArtists()
		return dto.NewArtistCollection(artists)
	}
}
