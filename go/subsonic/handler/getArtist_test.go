package handler

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func TestGetMissingArtist(t *testing.T) {

	db := dal.NewMock()
	handler := NewGetArtistHandler(db)
	params := url.Values{}
	url.Values.Add(params, "id", "3")

	context := api.HandlerContext{}

	response := handler(params, &context)

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

	if r.Message != "Artist not found." {
		t.Error("Wrong error message")
	}
}

func TestGetGoodArtist(t *testing.T) {

	db := dal.NewMock()
	handler := NewGetArtistHandler(db)
	params := url.Values{}
	url.Values.Add(params, "id", "1")

	context := api.HandlerContext{}

	response := handler(params, &context)

	if !response.IsSuccess {
		t.Error("Not a success")
	}

	r, ok := response.Body.(*dto.ArtistWithAlbums)
	if !ok {
		t.Error("Not an artist")
	}

	if r.ID != 1 {
		t.Error("Wrong ID")
	}

	if r.AlbumCount == 0 {
		t.Error("Missing albums")
	}
}

func TestGetArtistWithBadId(t *testing.T) {

	db := dal.NewMock()
	handler := NewGetArtistHandler(db)
	params := url.Values{}
	url.Values.Add(params, "id", "dfggsgs")

	context := api.HandlerContext{}

	response := handler(params, &context)

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

func TestGetArtistWithNoId(t *testing.T) {

	db := dal.NewMock()
	handler := NewGetArtistHandler(db)
	params := url.Values{}

	context := api.HandlerContext{}

	response := handler(params, &context)

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
