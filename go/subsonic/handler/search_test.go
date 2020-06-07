package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/subsonic/handler"
)

func TestSongSearch(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewSearchHandler(dal, handler.Search2)
	params := url.Values{}
	params.Add("query", "TITLE_2")

	response := handler(params, &api.HandlerContext{})
	if !response.IsSuccess {
		t.Error("Handler failed")
	}

	r, ok := response.Body.(*dto.Search2Response)
	if !ok {
		t.Error("Not a search response")
	}

	if len(r.Songs) != 5 {
		t.Error("Wrong song count")
	}
}

func TestArtistSearch(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewSearchHandler(dal, handler.Search2)
	params := url.Values{}
	params.Add("query", "artist_")

	response := handler(params, &api.HandlerContext{})
	if !response.IsSuccess {
		t.Error("Handler failed")
	}

	r, ok := response.Body.(*dto.Search2Response)
	if !ok {
		t.Error("Not a search response")
	}

	if len(r.Artists) != 5 {
		t.Error("Wrong artist count")
	}
}

func TestAlbumSearch(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewSearchHandler(dal, handler.Search2)
	params := url.Values{}
	params.Add("query", "album_1")

	response := handler(params, &api.HandlerContext{})
	if !response.IsSuccess {
		t.Error("Handler failed")
	}

	r, ok := response.Body.(*dto.Search2Response)
	if !ok {
		t.Error("Not a search response")
	}

	if len(r.Albums) != 1 {
		t.Error("Wrong album count")
	}
}
