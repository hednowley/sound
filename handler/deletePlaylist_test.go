package handler

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/dto"
)

func TestDeletePlaylist(t *testing.T) {

	db := dao.NewMockDatabase()

	handler := NewDeletePlaylistHandler(db)
	params := url.Values{}
	params.Add("id", "1")

	response := handler(params)

	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	if response.Body != nil {
		t.Error("Should have no body")
	}
}

func TestDeleteMissingPlaylist(t *testing.T) {

	db := dao.NewMockDatabase()

	handler := NewDeletePlaylistHandler(db)
	params := url.Values{}
	params.Add("id", "63")

	response := handler(params)

	if response.IsSuccess {
		t.Error("Should fail")
	}

	r, ok := response.Body.(dto.Error)
	if !ok {
		t.Error("Should return an error")
	}

	if r.Code != int(dto.NotFound) || r.Message != "Playlist not found: 63" {
		t.Error("Wrong error")
	}
}

func TestDeleteNonsensePlaylist(t *testing.T) {

	db := dao.NewMockDatabase()

	handler := NewDeletePlaylistHandler(db)
	params := url.Values{}
	params.Add("id", "shsd")

	response := handler(params)

	if response.IsSuccess {
		t.Error("Should fail")
	}

	r, ok := response.Body.(dto.Error)
	if !ok {
		t.Error("Should return an error")
	}

	if r.Code != int(dto.MissingParameter) || r.Message != "Required param (id) is missing" {
		t.Error("Wrong error")
	}
}

func TestDeletePlaylistNoParams(t *testing.T) {

	db := dao.NewMockDatabase()

	handler := NewDeletePlaylistHandler(db)
	params := url.Values{}

	response := handler(params)

	if response.IsSuccess {
		t.Error("Should fail")
	}

	e, ok := response.Body.(dto.Error)
	if !ok {
		t.Error("Should return an error")
	}

	if e.Code != int(dto.MissingParameter) || e.Message != "Required param (id) is missing" {
		t.Error("Wrong error")
	}
}
