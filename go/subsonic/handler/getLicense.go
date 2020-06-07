package handler

import (
	"net/url"

	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetLicenseHandler() api.Handler {
	return func(params url.Values, _ *api.HandlerContext) *api.Response {
		return api.NewSuccessfulReponse(dto.NewLicense())
	}
}
