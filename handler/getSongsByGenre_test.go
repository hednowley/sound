package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
	"github.com/hednowley/sound/handler"
)

func NewGetSongsByGenreTestResource() HandlerTestResource {
	db := dao.NewMockDatabase()
	handler := handler.NewGetSongsByGenreHandler(db)
	params := url.Values{}

	return HandlerTestResource{
		db:      db,
		handler: handler,
		params:  params,
	}
}

func TestGetSongsOfMissingGenre(t *testing.T) {

	h := NewGetSongsByGenreTestResource()
	url.Values.Add(h.params, "genre", "ss6s7 ssd")
	response := h.handler(h.params)

	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	r, ok := response.Body.(*dto.SongsByGenre)
	if !ok {
		t.Error("Wrong body")
	}

	if len(r.Songs) != 0 {
		t.Error("Shouldn't have songs")
	}
}

func TestGetSongsByGenre(t *testing.T) {

	h := NewGetSongsByGenreTestResource()
	url.Values.Add(h.params, "genre", "genre_1")
	response := h.handler(h.params)

	if !response.IsSuccess {
		t.Error("Not a success")
	}

	r, ok := response.Body.(*dto.SongsByGenre)
	if !ok {
		t.Error("Wrong body")
	}

	if len(r.Songs) != 2 || r.Songs[0].ID != 1 || r.Songs[1].ID != 10 {
		t.Error("Wrong songs")
	}
}

func TestGetSongsWithNoGenre(t *testing.T) {

	h := NewGetSongsByGenreTestResource()
	response := h.handler(h.params)

	if response.IsSuccess {
		t.Error("Not a failure")
	}

	r, ok := response.Body.(dto.Error)
	if !ok {
		t.Error("Not an error")
	}

	if r.Code != int(dto.MissingParameter) {
		t.Error("Wrong error code")
	}

	if r.Message != "Required param (genre) is missing" {
		t.Error("Wrong error message")
	}
}
