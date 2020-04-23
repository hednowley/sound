package handler

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/dto"
)

func TestGetArtists(t *testing.T) {

	db := dal.NewMock()
	handler := NewGetArtistsHandler(db, &config.Config{})
	params := url.Values{}
	url.Values.Add(params, "id", "1")
	response := handler(params)

	if !response.IsSuccess {
		t.Error("Not a success")
	}

	_, ok := response.Body.(*dto.ArtistCollection)
	if !ok {
		t.Error("Not an artist collection")
	}
}
