package handlers

import (
	"encoding/json"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/ws"
	"github.com/hednowley/sound/ws/dto"
)

func MakeGetArtistHandler(dal *dal.DAL) ws.WsHandler {
	return func(request *dto.Request) interface{} {
		var id uint

		if request.Params["id"] == nil || json.Unmarshal(*request.Params["id"], &id) != nil {
			return "bad id"
		}

		artist, err := dal.GetArtist(id)
		if err != nil {
			return "no artist"
		}
		return dto.NewArtist(artist)
	}
}
