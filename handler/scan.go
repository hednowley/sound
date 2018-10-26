package handler

import (
	"net/url"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dto"
)

func NewStartScanHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		go dal.StartAllScans(false, false)
		r := dto.NewScanStatus(dal.GetScanStatus(), dal.GetScanFileCount())
		return api.NewSuccessfulReponse(r)
	}
}

func NewGetScanStatusHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {
		r := dto.NewScanStatus(dal.GetScanStatus(), dal.GetScanFileCount())
		return api.NewSuccessfulReponse(r)
	}
}
