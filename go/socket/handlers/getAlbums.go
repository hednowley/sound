package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetAlbumsHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request) interface{} {
		albums, err := dal.Db.GetAlbums(dao.AlphabeticalByName, 9999999, 0)
		if err != nil {
			return dto.NewErrorResponse("error", 0)
		}
		return dto.NewAlbumCollection(albums)
	}
}
