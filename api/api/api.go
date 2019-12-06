// Package api is stuff.
package api

import (
	"net/http"

	"github.com/hednowley/sound/config"
)

type ControllerContext struct {
	// Pointer to a DTO struct. This struct is kept in a closure along
	// with the Run func to make Run like a generic function.
	// The struct must be mutated rather than reassigning the pointer.
	Body interface{}

	// Run the controller action.
	Run func(user *config.User, w http.ResponseWriter, r *http.Request) *Response
}

// Controller is a web controller.
// It accepts a data-transfer object and returns an unserialised Response.
type Controller struct {

	// Request token will be authenticated iff this is true
	Secure bool

	Make func() *ControllerContext
}

// BinaryController is a low-level web controller.
// It accepts a data-transfer object, a ResponseWriter and a Request.
// It returns a nil pointer to indicate that no further response is needed,
// otherwise it returns an unserialised Response.
type BinaryController struct {
	Secure bool                                                             // Request token will be authenticated iff this is true
	Run    func(http.ResponseWriter, *http.Request, *config.User) *Response // Run the controller action.
}
