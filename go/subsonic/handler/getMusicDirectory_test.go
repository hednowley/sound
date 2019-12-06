package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/handler"
)

func TestGetArtistDirectory(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewGetMusicDirectoryHandler(dal)
	params := url.Values{}
	url.Values.Add(params, "id", "artist_1")
	handler(params)

	// TODO
}
