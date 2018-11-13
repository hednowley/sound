package api

import (
	"net/http"
	"net/url"
)

var version = "1.15.0"

// Handler is a web controller action.
// It accepts a set of parameters and returns an unserialised Response.
type Handler func(url.Values) *Response

// BinaryHandler is a low-level web controller action.
// It accepts a set of parameters, a ResponseWriter and a Request.
// It returns a nil pointer to indicate that no further response is needed,
// otherwise it returns an unserialised Response.
type BinaryHandler func(url.Values, *http.ResponseWriter, *http.Request) *Response
