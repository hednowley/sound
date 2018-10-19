package handler

import (
	"net/url"

	log "github.com/cihub/seelog"
	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dto"

	"github.com/hednowley/sound/services"
)

func NewStartScanHandler(scanner *services.Scanner) api.Handler {

	return func(params url.Values) *api.Response {
		log.Info("Starting scan.")
		scanner.Start()

		r := dto.NewScanStatus(scanner.InProgress, scanner.FileCount())
		return api.NewSuccessfulReponse(r)
	}
}

func NewGetScanStatusHandler(scanner *services.Scanner) api.Handler {

	return func(params url.Values) *api.Response {
		r := dto.NewScanStatus(scanner.InProgress, scanner.FileCount())
		return api.NewSuccessfulReponse(r)
	}
}
