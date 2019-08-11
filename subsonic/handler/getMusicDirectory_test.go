package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/subsonic/handler"
)

func TestGetArtistDirectory(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewGetMusicDirectoryHandler(dal)
	params := url.Values{}
	url.Values.Add(params, "id", "artist_1")
	response := handler(params)

	if response.IsSuccess {
		t.Error("Not a failure")
	}

	r, ok := response.Body.(dto.Error)
	if !ok {
		t.Error("Not an error")
	}

	if r.Code != int(dto.NotFound) {
		t.Error("Wrong error code")
	}

	if r.Message != "Album not found." {
		t.Error("Wrong error message")
	}
}
