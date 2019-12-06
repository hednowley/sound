package api

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestEmptySerialiser(t *testing.T) {

	r := Response{
		Body:      nil,
		IsSuccess: true,
	}

	j := serialiseToJSON(&r)
	if j != fmt.Sprintf(`{"subsonic-response":{"status":"ok","version":"%v"}}`, version) {
		t.Error()
	}

	r.IsSuccess = false

	x := serialiseToXML(&r)
	if x != fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="%v"></subsonic-response>`, version) {
		t.Error()
	}
}

func TestSerialiser(t *testing.T) {

	type TestResponseDto struct {
		XMLName xml.Name `xml:"test-thing" json:"-"`
		ID      uint     `xml:"id,attr" json:"id,string"`
		Thingy  string   `xml:"thingy,attr" json:"thingy"`
	}

	r := Response{
		Body: &TestResponseDto{
			ID:     5,
			Thingy: "dyig",
		},
		IsSuccess: false,
	}

	j := serialiseToJSON(&r)
	if j != fmt.Sprintf(`{"subsonic-response":{"status":"failed","version":"%v", "test-thing": {"id":"5","thingy":"dyig"}}}`, version) {
		t.Error()
	}

	r.IsSuccess = true

	x := serialiseToXML(&r)
	if x != fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="%v"><test-thing id="5" thingy="dyig"></test-thing></subsonic-response>`, version) {
		t.Error()
	}
}
