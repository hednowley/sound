package handler

import (
	"net/url"

	"github.com/hednowley/sound/idal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewStartScanHandler(dal idal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		go dal.StartAllScans(false, false)
		r := dto.NewScanStatus(dal.GetScanStatus(), dal.GetScanFileCount())
		return api.NewSuccessfulReponse(r)
	}
}

func NewGetScanStatusHandler(dal idal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		r := dto.NewScanStatus(dal.GetScanStatus(), dal.GetScanFileCount())
		return api.NewSuccessfulReponse(r)
	}
}
