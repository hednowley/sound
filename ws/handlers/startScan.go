package handlers

import (
	"encoding/json"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/ws/dto"
)

func MakeStartScanHandler(dal *dal.DAL) WsHandler {
	return func(request *dto.Request) interface{} {
		var update bool
		err := json.Unmarshal(*request.Params["update"], &update)
		if err != nil {

		}

		go dal.StartAllScans(update, false)
		return struct{}{}
	}
}
