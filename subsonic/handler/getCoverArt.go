package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cihub/seelog"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/services"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetCoverArtHandler(dal interfaces.DAL) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request) *api.Response {

		id := params.Get("id")
		if id == "" {
			return api.NewErrorReponse(dto.NotFound, "Art parameter missing")
		}

		path, err := dal.GetArtPath(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		sizeParam := params.Get("size")
		size := api.ParseUint(sizeParam, 0)
		if size == 0 {
			http.ServeFile(*w, r, path)
			return nil
		}

		dir := filepath.Dir(path)
		ext := filepath.Ext(path)
		resized := filepath.Join(dir, fmt.Sprintf("%v_%v%v", art.ID, size, ext))
		_, err = os.Stat(resized)
		if os.IsNotExist(err) {
			seelog.Tracef("Resizing %v to %v", id, size)
			err = services.Resize(path, resized, uint(size))
			if err != nil {
				seelog.Errorf("Could not resize artwork %v", id)
				resized = path
			}
		}

		http.ServeFile(*w, r, resized)
		return nil
	}
}
