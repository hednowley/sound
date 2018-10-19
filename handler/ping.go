package handler

import (
	"net/url"

	"github.com/hednowley/sound/api"
)

// NewPingHandler is a handler for responding to ping requests.
// It replies to any request with an empty success response.
func NewPingHandler() api.Handler {
	return func(params url.Values) *api.Response {
		return &api.Response{
			Body:      nil,
			IsSuccess: true,
		}
	}
}
