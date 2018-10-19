package handler

import (
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {

	handler := NewPingHandler()
	response := handler(url.Values{})

	if !response.IsSuccess {
		t.Error()
	}

	if response.Body != nil {
		t.Error()
	}
}
