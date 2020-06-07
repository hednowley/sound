package handler_test

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/handler"
)

func TestUpdatePlaylist(t *testing.T) {

	db := dal.NewMock()

	handler := handler.NewUpdatePlaylistHandler(db)
	params := url.Values{}
	params.Add("playlistId", "1")
	params.Add("name", "new_name")

	context := api.HandlerContext{
		User: &config.User{
			Username: "tommy",
		},
	}

	response := handler(params, &context)

	if !response.IsSuccess {
		t.Error("Should succeed")
	}

	if response.Body != nil {
		t.Error("Should have no body")
	}
}
