package handlers

import (
	"encoding/json"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeGetPlaylistHandler(dal *dal.DAL) socket.Handler {
	return func(request *dto.Request) interface{} {
		var id uint

		if request.Params["id"] == nil || json.Unmarshal(*request.Params["id"], &id) != nil {
			return "bad id"
		}

		playlist, err := dal.GetPlaylist(id)
		if err != nil {
			return "no playlist"
		}
		return dto.NewPlaylist(playlist)
	}
}
