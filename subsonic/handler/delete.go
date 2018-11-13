package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
)

func NewDeleteHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		dal.Empty()
		return &api.Response{
			Body:      nil,
			IsSuccess: true,
		}
	}
}
