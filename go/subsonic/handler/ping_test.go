package handler

import (
	"net/url"
	"testing"

	"github.com/hednowley/sound/subsonic/api"
)

func TestPing(t *testing.T) {

	handler := NewPingHandler()
	response := handler(url.Values{}, &api.HandlerContext{})

	if !response.IsSuccess {
		t.Error()
	}

	if response.Body != nil {
		t.Error()
	}
}
