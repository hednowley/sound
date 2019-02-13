package controller

import (
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
)

func NewStartScanHandler(dal *dal.DAL) *api.Controller {

	input := struct {
		Update bool
		Delete bool
	}{}

	w := func(_ *config.User) *api.Response {
		go dal.StartAllScans(input.Update, input.Delete)
		r := dto.NewScanStatus(dal.GetScanStatus(), dal.GetScanFileCount())
		return api.NewSuccessfulReponse(&r)
	}

	return &api.Controller{
		Input:  &input,
		Run:    w,
		Secure: true,
	}
}

func NewGetScanStatusHandler(dal *dal.DAL) *api.Controller {

	input := struct{}{}

	w := func(_ *config.User) *api.Response {
		r := dto.NewScanStatus(dal.GetScanStatus(), dal.GetScanFileCount())
		return api.NewSuccessfulReponse(&r)
	}

	return &api.Controller{
		Input:  &input,
		Run:    w,
		Secure: true,
	}
}
