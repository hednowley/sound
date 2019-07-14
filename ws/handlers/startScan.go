package handlers

import (
	"encoding/json"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/ws/dto"
)

func MakeStartScanHandler(scanner *provider.Scanner) interfaces.WsHandler {
	return func(request *dto.Request) interface{} {
		var update bool
		var delete bool

		err := json.Unmarshal(*request.Params["update"], &update)
		if err != nil {

		}

		err = json.Unmarshal(*request.Params["delete"], &delete)
		if err != nil {

		}

		go scanner.StartAllScans(update, delete)
		return struct{}{}
	}
}
