package handler

import (
	"net/http"
	"net/url"
	"os"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewGetCoverArtHandler(dal *dal.DAL) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request) *api.Response {

		id := params.Get("id")
		if id == "" {
			return api.NewErrorReponse(dto.NotFound, "Art parameter missing")
		}

		sizeParam := params.Get("size")
		size := util.ParseUint(sizeParam, 0)

		path := dal.GetArtPath(id, size)

		_, err := os.Stat(path)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		http.ServeFile(*w, r, path)
		return nil
	}
}
