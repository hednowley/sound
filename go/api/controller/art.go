package controller

import (
	"net/http"
	"os"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/util"
)

// NewArtController creates a controller for serving artwork images.
func NewArtController(dal *dal.DAL) *api.BinaryController {

	run := func(w http.ResponseWriter, r *http.Request, _ *config.User) *api.Response {

		params := r.URL.Query()
		id := params.Get("id")

		sizeParam := params.Get("size")
		size := util.ParseUint(sizeParam, 0)

		path := dal.GetArtPath(id, size)

		_, err := os.Stat(path)
		if err != nil {
			return api.NewErrorReponse(err.Error())
		}

		http.ServeFile(w, r, path)
		return nil
	}

	return &api.BinaryController{
		Run:    run,
		Secure: true,
	}
}
