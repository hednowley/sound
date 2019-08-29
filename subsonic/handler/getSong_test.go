package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/subsonic/handler"
)

type HandlerTestResource struct {
	db      *dal.DAL
	handler api.Handler
	params  url.Values
}

func NewGetSongTestResource() HandlerTestResource {
	db := dal.NewMock()
	handler := handler.NewGetSongHandler(db)
	params := url.Values{}

	return HandlerTestResource{
		db:      db,
		handler: handler,
		params:  params,
	}
}

func TestGetMissingSong(t *testing.T) {

	h := NewGetSongTestResource()
	url.Values.Add(h.params, "id", "666")
	response := h.handler(h.params)

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

	if r.Message != "Song not found." {
		t.Error("Wrong error message")
	}
}

func TestGetGoodSong(t *testing.T) {

	h := NewGetSongTestResource()
	url.Values.Add(h.params, "id", "1")
	response := h.handler(h.params)

	if !response.IsSuccess {
		t.Error("Not a success")
	}

	r, ok := response.Body.(*dto.Song)
	if !ok {
		t.Error("Not a song")
	}

	if r.ID != 1 {
		t.Error("Wrong ID")
	}
}

func TestGetSongWithBadId(t *testing.T) {

	h := NewGetSongTestResource()
	url.Values.Add(h.params, "id", "dfggsgs")
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

	if r.Message != "Required param (id) is missing" {
		t.Error("Wrong error message")
	}
}

func TestGetSongWithNoId(t *testing.T) {

	h := NewGetSongTestResource()
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

	if r.Message != "Required param (id) is missing" {
		t.Error("Wrong error message")
	}
}
