package handler

import (
	"net/http"
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

func NewGetCoverArtHandler(dal *dal.DAL) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request, _ *api.HandlerContext) *api.Response {

		id := params.Get("id")
		if id == "" {
			return api.NewErrorReponse(dto.NotFound, "Art parameter missing")
		}

		sizeParam := params.Get("size")
		size := util.ParseUint(sizeParam, 0)

		path := dal.GetArt(id, size)
		if path == nil {
			return api.NewErrorReponse(dto.Generic, "Art is not available")
		}

		http.ServeFile(*w, r, *path)
		return nil
	}
}
