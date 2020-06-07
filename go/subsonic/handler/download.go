package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewDownloadHandler does http://www.subsonic.org/pages/api.jsp#download
func NewDownloadHandler(dal *dal.DAL) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request, _ *api.HandlerContext) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.Generic, fmt.Sprintf("Song not found: %v", idParam))
		}

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		path := dal.Db.GetSongPath(conn, id)
		if path != nil {
			return api.NewErrorReponse(dto.Generic, "song not found")
		}

		http.ServeFile(*w, r, *path)
		return nil
	}
}
