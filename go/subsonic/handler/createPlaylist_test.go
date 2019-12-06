package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/subsonic/handler"
)

func TestCreatePlaylist(t *testing.T) {

	db := dal.NewMock()

	handler := handler.NewCreatePlaylistHandler(db)
	params := url.Values{}
	params.Add("name", "dfsdf")
	params.Add("songId", "1")
	params.Add("songId", "-50")
	params.Add("songId", "15")
	params.Add("songId", "10")
	params.Add("songId", "dfg")

	response := handler(params)

	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	r, ok := response.Body.(*dto.Playlist)
	if !ok {
		t.Error("Should return a playlist")
	}

	if r.Name != "dfsdf" {
		t.Error("Wrong name")
	}

	if len(r.Songs) != 3 || r.Songs[0].ID != 1 || r.Songs[1].ID != 15 || r.Songs[2].ID != 10 {
		t.Error("Wrong songs")
	}
}

func TestEmptyPlaylist(t *testing.T) {

	db := dal.NewMock()

	handler := handler.NewCreatePlaylistHandler(db)
	params := url.Values{}
	params.Add("playlistId", "1")
	params.Add("songId", "sdf")
	params.Add("songId", "3")

	response := handler(params)

	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	r, ok := response.Body.(*dto.Playlist)
	if !ok {
		t.Error("Should return a playlist")
	}

	if r.Name != "playlist_1" {
		t.Error("Wrong name")
	}

	if len(r.Songs) != 0 {
		t.Error("Shouldn't have any songs")
	}
}

func TestReplacePlaylist(t *testing.T) {

	db := dal.NewMock()

	handler := handler.NewCreatePlaylistHandler(db)
	params := url.Values{}
	params.Add("playlistId", "1")
	params.Add("name", "dsfsad")
	params.Add("songId", "sdf")
	params.Add("songId", "14")
	params.Add("songId", "10")

	response := handler(params)

	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	r, ok := response.Body.(*dto.Playlist)
	if !ok {
		t.Error("Should return a playlist")
	}

	if r.Name != "dsfsad" {
		t.Error("Wrong name")
	}

	if len(r.Songs) != 2 || r.Songs[0].ID != 14 || r.Songs[1].ID != 10 {
		t.Error("Wrong songs")
	}
}

func TestReplaceMissingPlaylist(t *testing.T) {

	db := dal.NewMock()

	handler := handler.NewCreatePlaylistHandler(db)
	params := url.Values{}
	params.Add("playlistId", "666")
	params.Add("name", "dsfsad")
	params.Add("songId", "sdf")
	params.Add("songId", "14")
	params.Add("songId", "10")

	response := handler(params)

	if response.IsSuccess {
		t.Error("Should fail")
	}

	e, ok := response.Body.(dto.Error)
	if !ok {
		t.Error("Should return an error")
	}

	if e.Code != int(dto.NotFound) || e.Message != "Playlist not found: 666" {
		t.Error("Wrong error")
	}
}

func TestCreatePlaylistNoParams(t *testing.T) {

	db := dal.NewMock()

	handler := handler.NewCreatePlaylistHandler(db)
	params := url.Values{}

	response := handler(params)

	if response.IsSuccess {
		t.Error("Should fail")
	}

	e, ok := response.Body.(dto.Error)
	if !ok {
		t.Error("Should return an error")
	}

	if e.Code != int(dto.MissingParameter) || e.Message != "Playlist ID or name must be specified." {
		t.Error("Wrong error")
	}
}
