package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
)

func getStatus(succeeded bool) string {
	if succeeded {
		return "ok"
	}
	return "failed"
}

func serialiseToJSON(response *Response) string {
	var body *string
	var name string

	if response.Body != nil {
		b, err := json.Marshal(response.Body)
		if err == nil {
			s := string(b)
			body = &s

			// Dereference response.Body if it is a pointer
			v := reflect.ValueOf(response.Body)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			nameField, _ := v.Type().FieldByName("XMLName")
			name = nameField.Tag.Get("xml")
		}
	}

	responseJSON, _ := json.Marshal(response.Body)

	status := getStatus(response.IsSuccess)

	if body == nil || len(name) == 0 {
		return fmt.Sprintf(`{"subsonic-response":{"status":"%v","version":"%v"}}`,
			status, version)
	}

	return fmt.Sprintf(
		`{"subsonic-response":{"status":"%v","version":"%v", "%v": %s}}`,
		status, version, name, responseJSON)
}

func serialiseToXML(response *Response) string {
	var body string
	if response.Body != nil {
		b, err := xml.Marshal(response.Body)
		if err == nil {
			body = string(b)
		}
	}

	return fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8"?><subsonic-response xmlns="http://subsonic.org/restapi" status="%v" version="%v">%v</subsonic-response>`,
		getStatus(response.IsSuccess),
		version,
		body)
}
