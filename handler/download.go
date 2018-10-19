package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

// NewDownloadHandler is a handler for downloading songs.
func NewDownloadHandler(database *dao.Database) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request) *api.Response {

		idParam := params.Get("id")
		id := api.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.Generic, fmt.Sprintf("Song not found: %v", idParam))
		}

		file, err := database.GetSong(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		http.ServeFile(*w, r, file.Path)
		return nil
	}
}
