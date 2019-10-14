package controller

import (
	"net/http"
	"path"
	"strings"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/util"
)

func NewStreamController(dal interfaces.DAL) *api.BinaryController {

	run := func(w *http.ResponseWriter, r *http.Request, _ *config.User) *api.Response {

		params := r.URL.Query()
		idStr := params.Get("id")
		id := util.ParseUint(idStr, 0)

		if id == 0 {
			return api.NewErrorReponse("No ID!")
		}

		file, err := dal.GetSong(id, false, false, false)
		if err != nil {
			return api.NewErrorReponse(err.Error())
		}

		// http.ServeFile incorrectly guesses the Content-Type as "video/mp4" for AAC files
		// so we override it here.
		ext := strings.ToLower(path.Ext(file.Path))
		if ext == ".aac" || ext == ".m4a" {
			(*w).Header()["Content-Type"] = []string{"audio/aac"}
		}

		http.ServeFile(*w, r, file.Path)
		return nil
	}

	return &api.BinaryController{
		Run:    run,
		Secure: true,
	}
}
