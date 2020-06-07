package api

import (
	"net/http"
	"net/url"

	"github.com/hednowley/sound/config"
)

var version = "1.16.1"

type HandlerContext struct {
	User *config.User
}

// Handler is a web controller action.
// It accepts a set of parameters and returns an unserialised Response.
type Handler func(url.Values, *HandlerContext) *Response

// BinaryHandler is a low-level web controller action.
// It accepts a set of parameters, a ResponseWriter and a Request.
// It returns a nil pointer to indicate that no response is needed,
// otherwise it returns an unserialised Response.
type BinaryHandler func(url.Values, *http.ResponseWriter, *http.Request, *HandlerContext) *Response
