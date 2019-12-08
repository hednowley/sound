package handlers

import (
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/ws"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetAlbumsHandler(dal *dal.DAL) ws.WsHandler {
	return func(request *dto.Request) interface{} {
		albums := dal.GetAlbums(dao.AlphabeticalByName, 9999999, 0)
		return dto.NewAlbumCollection(albums)
	}
}
