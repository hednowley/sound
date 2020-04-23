package handlers

import (
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetArtistsHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request) interface{} {
		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse("Cannot make DB conn")
		}
		defer conn.Release()

		artists, err := dal.Db.GetArtists(conn)
		if err != nil {
			return dto.NewErrorResponse("error", 0)
		}
		return dto.NewArtistCollection(artists)
	}
}
