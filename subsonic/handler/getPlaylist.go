package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetPlaylistHandler(dal interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("id")
		id := api.ParseUint(idParam, 0)
		if id == 0 {
			message := fmt.Sprintf("Playlist not found: %v", idParam)
			return api.NewErrorReponse(dto.Generic, message)
		}

		p, err := dal.GetPlaylist(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewPlaylist(p))
	}
}
