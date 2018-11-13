package api

import (
	"net/url"
	"strconv"
	"strings"
)

// parseResponseFormat tries to parse a string into a response format.
// If this is not possible then a nil pointer is returned instead.
func parseResponseFormat(param string) *responseFormat {
	if len(param) == 0 {
		return &defaultFormat
	}

	param = strings.ToLower(param)
	if param == "json" {
		f := jsonFormat
		return &f
	}

	if param == "xml" {
		f := xmlFormat
		return &f
	}

	return nil
}

// ParseUint tries to parse a string into a uint. If this is not possible
// then the given default is returned instead.
func ParseUint(param string, defaultValue uint) uint {

	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return defaultValue
	}
	return uint(id)
}

// ParseBool tries to parse a string into a bool. If this is not possible
// then a nil pointer is returned instead.
func ParseBool(param string) *bool {

	var b bool
	param = strings.ToLower(param)
	if param == "true" || param == "1" {
		b = true
		return &b
	}
	if param == "false" || param == "0" {
		b = false
		return &b
	}

	return nil
}

// parseParams extracts parameters from a Request.
// It merges URL queries and the request body
// (body wins where there are collisions)
func parseParams(urlQuery url.Values, body []byte) url.Values {

	values := url.Values{}

	for k, v := range urlQuery {
		for _, s := range v {
			values.Add(k, s)
		}
	}

	bodyParams, err := url.ParseQuery(string(body))
	if err == nil {
		for k, v := range bodyParams {
			for _, s := range v {
				values.Add(k, s)
			}
		}
	}

	return values
}
