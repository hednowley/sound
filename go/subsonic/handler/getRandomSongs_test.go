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

	context := api.HandlerContext{}

	response := handler(params, &context)

	songs := testResponse(response, t)

	if len(songs.Songs) != 2 {
		t.Error("Should have two songs")
		return
	}

	s := songs.Songs[0]

	if s.ID == songs.Songs[1].ID {
		t.Error("Songs should be different")
	}
}

func TestGetRandomGenreSongs(t *testing.T) {

	dal := dal.NewMock()
	handler := handler.NewGetRandomSongsHandler(dal)
	params := url.Values{}
	url.Values.Add(params, "genre", "geNRe_1")

	context := api.HandlerContext{}

	response := handler(params, &context)

	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	songs := testResponse(response, t)

	if len(songs.Songs) != 3 {
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

	context := api.HandlerContext{}

	response := handler(params, &context)
	songs := testResponse(response, t)
	if len(songs.Songs) != 8 {
		t.Error("Nonsense years should return all songs")
	}

	params = url.Values{}
	url.Values.Add(params, "fromYear", "2000")
	url.Values.Add(params, "toYear", "2010")

	response = handler(params, &context)
	songs = testResponse(response, t)
	if len(songs.Songs) != 3 {
		t.Error()
	}

	url.Values.Set(params, "toYear", "2008")

	response = handler(params, &context)
	songs = testResponse(response, t)
	if len(songs.Songs) != 2 {
		t.Error()
	}
}
