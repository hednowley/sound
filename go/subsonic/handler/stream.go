package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
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

// NewStreamHandler does http://www.subsonic.org/pages/api.jsp#stream
func NewStreamHandler(dal *dal.DAL) api.BinaryHandler {

	return func(params url.Values, w *http.ResponseWriter, r *http.Request) *api.Response {

		idParam := params.Get("id")
		id := util.ParseUint(idParam, 0)
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

		file, err := dal.GetSong(id)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		// http.ServeFile incorrectly guesses the Content-Type as "video/mp4" for AAC files
		// so we override it here.
		ext := strings.ToLower(path.Ext(file.Path))
		if ext == ".aac" || ext == ".m4a" {
			(*w).Header()["Content-Type"] = []string{"audio/aac"}
		}

		http.ServeFile(*w, r, file.Path)
		return nil
	}
}
