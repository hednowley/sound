package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewDownloadHandler is a handler for downloading songs.
func NewDownloadHandler(database interfaces.DAL) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.Generic, fmt.Sprintf("Song not found: %v", idParam))
		}

		file, err := database.GetSong(id, false, false, false)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		http.ServeFile(*w, r, file.Path)
		return nil
	}
}
