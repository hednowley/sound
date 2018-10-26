package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dto"
)

func NewGetCoverArtHandler(dal *dal.DAL) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request) *api.Response {

		idParam := params.Get("id")
		id := api.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.NotFound, fmt.Sprintf("Art not found: %v", idParam))
		}

		//sizeParam := params.Get("size")
		//size := api.ParseUint(sizeParam, 0)

		art, err := dal.GetArt(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		//if size == 0 {
		http.ServeFile(*w, r, art.Path)
		return nil
		//}

		/*
			dir := filepath.Dir(art.Path)
			ext := filepath.Ext(art.Path)
			resized := filepath.Join(dir, fmt.Sprintf("%v_%v.%v", art.ID, size, ext))
			_, err = os.Stat(resized)
			if os.IsNotExist(err) {
				err = services.Resize(art.Path, resized, uint(size))
				if err != nil {
					resized = art.Path
				}
			}

			http.ServeFile(*w, r, resized)
			return nil
		*/
	}
}
