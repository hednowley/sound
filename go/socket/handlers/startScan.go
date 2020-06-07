package handlers

import (
	"encoding/json"

	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/dto"
)

func MakeStartScanHandler(scanner *provider.Scanner) socket.Handler {
	return func(request *dto.Request, _ *socket.HandlerContext) interface{} {
		var update bool
		var delete bool

		if request.Params["update"] == nil || json.Unmarshal(*request.Params["update"], &update) != nil {
			update = false
		}

		if request.Params["delete"] == nil || json.Unmarshal(*request.Params["delete"], &delete) != nil {
			delete = false
		}

		go scanner.StartAllScans(update, delete)
		return struct{}{}
	}
}
