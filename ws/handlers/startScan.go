package handlers

import (
	"encoding/json"

	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/ws/dto"
)

func MakeStartScanHandler(dal idal.DAL) WsHandler {
	return func(request *dto.Request) interface{} {
		var update bool
		err := json.Unmarshal(*request.Params["update"], &update)
		if err != nil {

		}

		go dal.StartAllScans(update, false)
		return struct{}{}
	}
}
