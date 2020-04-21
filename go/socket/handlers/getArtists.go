package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetArtistsHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request) interface{} {
		artists, err := dal.Db.GetArtists()
		if err != nil {
			return dto.NewErrorResponse("error", 0)
		}
		return dto.NewArtistCollection(artists)
	}
}
