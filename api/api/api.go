// Package api is stuff.
package api

import (
	"net/http"

	"github.com/hednowley/sound/config"
)

// Controller is a web controller.
// It accepts a data-transfer object and returns an unserialised Response.
type Controller struct {

	// Request token will be authenticated iff this is true
	Secure bool

	// Run the controller action.
	Run func(*config.User) *Response
}

// BinaryController is a low-level web controller.
// It accepts a data-transfer object, a ResponseWriter and a Request.
// It returns a nil pointer to indicate that no further response is needed,
// otherwise it returns an unserialised Response.
type BinaryController struct {
	Secure bool                                                              // Request token will be authenticated iff this is true
	Run    func(*http.ResponseWriter, *http.Request, *config.User) *Response // Run the controller action.
}
