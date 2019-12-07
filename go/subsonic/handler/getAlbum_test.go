package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/subsonic/handler"
)

func TestGetMissingAlbum(t *testing.T) {

	db := dal.NewMock()
	handler := handler.NewGetAlbumHandler(db)
	params := url.Values{}
	url.Values.Add(params, "id", "666")
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

func TestGetGoodAlbum(t *testing.T) {

	db := dal.NewMock()
	handler := handler.NewGetAlbumHandler(db)
	params := url.Values{}
	url.Values.Add(params, "id", "1")
	response := handler(params)

	if !response.IsSuccess {
		t.Error("Not a success")
	}

	r, ok := response.Body.(*dto.Album)
	if !ok {
		t.Error("Not an album")
	}

	if r.ID != 1 {
		t.Error("Wrong ID")
	}
}

func TestGetAlbumWithBadId(t *testing.T) {

	db := dal.NewMock()
	handler := handler.NewGetAlbumHandler(db)
	params := url.Values{}
	url.Values.Add(params, "id", "dfggsgs")
	response := handler(params)

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

func TestGetAlbumWithNoId(t *testing.T) {

	db := dal.NewMock()
	handler := handler.NewGetAlbumHandler(db)
	params := url.Values{}
	response := handler(params)

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