package controller

import (
	"net/http"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/interfaces"
)

func NewArtController(dal interfaces.DAL) *api.BinaryController {

	run := func(w *http.ResponseWriter, r *http.Request, _ *config.User) *api.Response {

		params := r.URL.Query()
		id := params.Get("id")

		path, err := dal.GetArtPath(id)
		if err != nil {
			return api.NewErrorReponse(err.Error())
		}

		http.ServeFile(*w, r, path)
		return nil
	}

	return &api.BinaryController{
		Run:    run,
		Secure: true,
	}
}
