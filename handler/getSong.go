package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

func NewGetSongHandler(database *dao.Database) api.Handler {

	return func(params url.Values) *api.Response {

		idParam := params.Get("id")
		id := api.ParseUint(idParam, 0)
		if id == 0 {
			message := fmt.Sprintf("Song not found: %v", idParam)
			return api.NewErrorReponse(dto.Generic, message)
		}

		file, err := database.GetSong(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		return api.NewSuccessfulReponse(dto.NewSong(file))
	}
}
