package handler

import (
	"net/url"

	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/subsonic/api"
)

func NewDeleteHandler(dal idal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		dal.Empty()
		return &api.Response{
			Body:      nil,
			IsSuccess: true,
		}
	}
}
