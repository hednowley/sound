package handler

import (
	"fmt"
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

// NewDeletePlaylistHandler does http://www.subsonic.org/pages/api.jsp#deletePlaylist
func NewDeletePlaylistHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values, context *api.HandlerContext) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.MissingParameter, "Required param (id) is missing")
		}

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		err = dal.Db.DeletePlaylist(conn, id, context.User.Username)
		if err != nil {
			if _, ok := err.(*dao.ErrNotFound); ok {
				message := fmt.Sprintf("Playlist not found: %v", idParam)
				return api.NewErrorReponse(dto.NotFound, message)
			}
			return api.NewErrorReponse(dto.NotFound, err.Error())
		}

		return api.NewSuccessfulReponse(nil)
	}
}
