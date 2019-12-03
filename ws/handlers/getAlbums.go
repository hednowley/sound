package handlers

import (
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetAlbumsHandler(dal interfaces.DAL) interfaces.WsHandler {
	return func(request *dto.Request) interface{} {
		albums := dal.GetAlbums(dao.AlphabeticalByName, 999999999999, 0)
		return dto.NewAlbumCollection(albums)
	}
}
