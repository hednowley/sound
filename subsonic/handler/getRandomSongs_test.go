package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/subsonic/handler"
)

func testResponse(response *api.Response, t *testing.T) *dto.RandomSongs {
	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	r, ok := response.Body.(*dto.RandomSongs)
	if !ok {
		t.Error("Wrong body")
	}

	return r
}

func TestGetRandomSongs(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewGetRandomSongsHandler(dal)
	params := url.Values{}
	url.Values.Add(params, "size", "2")

	response := handler(params)

	songs := testResponse(response, t)

	if len(songs.Songs) != 2 {
		t.Error("Should have two songs")
	}

	s := songs.Songs[0]

	if s.ID == songs.Songs[1].ID {
		t.Error("Songs should be different")
	}

	if s.AlbumName == "" || s.ArtistName == "" || s.Genre == "" {
		t.Error("Songs should have joins")
	}
}

func TestGetRandomGenreSongs(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewGetRandomSongsHandler(dal)
	params := url.Values{}
	url.Values.Add(params, "genre", "geNRe_1")

	response := handler(params)

	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	songs := testResponse(response, t)

	if len(songs.Songs) != 2 {
		t.Error("Should have two songs")
	}

	if songs.Songs[0].ID == songs.Songs[1].ID {
		t.Error("Songs should be different")
	}
}

func TestGetRandomSongsWithYears(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewGetRandomSongsHandler(dal)
	params := url.Values{}
	url.Values.Add(params, "fromYear", "asdiuhasdiuhasd")
	url.Values.Add(params, "toYear", "asdiuhasdiuhasd")

	response := handler(params)
	songs := testResponse(response, t)
	if len(songs.Songs) != 6 {
		t.Error()
	}

	params = url.Values{}
	url.Values.Add(params, "fromYear", "2000")
	url.Values.Add(params, "toYear", "2010")

	response = handler(params)
	songs = testResponse(response, t)
	if len(songs.Songs) != 3 {
		t.Error()
	}

	url.Values.Set(params, "toYear", "2008")

	response = handler(params)
	songs = testResponse(response, t)
	if len(songs.Songs) != 2 {
		t.Error()
	}
}
