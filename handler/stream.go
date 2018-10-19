package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

type StreamFormat int

const (
	MP3 StreamFormat = 0
	FLV StreamFormat = 1
	Raw StreamFormat = 2
)

func parseStreamFormat(param string) *StreamFormat {

	var f StreamFormat
	if len(param) == 0 {
		f = Raw
		return &f
	}

	param = strings.ToLower(param)
	if param == "mp3" {
		f = MP3
		return &f
	}

	if param == "flv" {
		f = FLV
		return &f
	}

	return nil
}

func NewStreamHandler(database *dao.Database) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request) *api.Response {

		idParam := params.Get("id")
		id := api.ParseUint(idParam, 0)
		if id == 0 {
			return api.NewErrorReponse(dto.Generic, fmt.Sprintf("Song not found: %v", idParam))
		}

		/*
			bitrateParam := params.Get("maxBitRate")
			bitrate, err := parseInt64(bitrateParam)
			if err != nil {
				return fmt.Errorf("Invalid bitrate: %v", bitrateParam)
			}

			formatParam := params.Get("format")
			format, err := parseStreamFormat(formatParam)
			if err != nil {
				return fmt.Errorf("Invalid format: %v", formatParam)
			}

			estimateParam := params.Get("estimateContentLength")
			estimate, err := parseBool(estimateParam)
			if err != nil {
				return fmt.Errorf("Invalid estimation flag: %v", estimateParam)
			}
		*/

		file, err := database.GetSong(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		http.ServeFile(*w, r, file.Path)
		return nil
	}
}
