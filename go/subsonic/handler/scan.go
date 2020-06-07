package handler

import (
	"net/url"

	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewStartScanHandler(dal *provider.Scanner) api.Handler {

	return func(params url.Values, _ *api.HandlerContext) *api.Response {
		go dal.StartAllScans(false, false)
		r := dto.NewScanStatus(dal.GetScanStatus(), dal.GetScanFileCount())
		return api.NewSuccessfulReponse(r)
	}
}

func NewGetScanStatusHandler(dal *provider.Scanner) api.Handler {

	return func(params url.Values, _ *api.HandlerContext) *api.Response {
		r := dto.NewScanStatus(dal.GetScanStatus(), dal.GetScanFileCount())
		return api.NewSuccessfulReponse(r)
	}
}
