package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cihub/seelog"
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/services"
	"github.com/hednowley/sound/util"
)

func NewArtController(dal *dal.DAL) *api.BinaryController {

	run := func(w http.ResponseWriter, r *http.Request, _ *config.User) *api.Response {

		params := r.URL.Query()
		id := params.Get("id")

		path, err := dal.GetArtPath(id)
		if err != nil {
			return api.NewErrorReponse(err.Error())
		}

		sizeParam := params.Get("size")
		size := util.ParseUint(sizeParam, 0)
		if size == 0 {
			http.ServeFile(w, r, path)
			return nil
		}

		dir, filename := filepath.Split(path)
		ext := filepath.Ext(filename)
		resized := filepath.Join(dir, fmt.Sprintf("%v_%v%v", strings.TrimSuffix(filename, ext), size, ext))
		_, err = os.Stat(resized)
		if os.IsNotExist(err) {
			seelog.Tracef("Resizing %v to %v", id, size)
			err = services.Resize(path, resized, uint(size))
			if err != nil {
				seelog.Errorf("Could not resize artwork %v", id)
				resized = path
			}
		}

		http.ServeFile(w, r, resized)
		return nil
	}

	return &api.BinaryController{
		Run:    run,
		Secure: true,
	}
}
