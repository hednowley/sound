// Package api is stuff.
package api

import (
	"net/http"

	"github.com/hednowley/sound/config"
)

// Controller is a web controller.
// It accepts a data-transfer object and returns an unserialised Response.
type Controller struct {
	// Pointer to a DTO struct. This struct should be kept in a closure along
	// with the Run func (makes Run like a generic function).
	// The struct must be mutated rather than reassigning the pointer.
	Input interface{}

	// Request token will be authenticated iff this is true
	Secure bool

	// Run the controller action.
	Run func(user *config.User) *Response
}

// BinaryController is a low-level web controller.
// It accepts a data-transfer object, a ResponseWriter and a Request.
// It returns a nil pointer to indicate that no further response is needed,
// otherwise it returns an unserialised Response.
type BinaryController struct {
	Input  interface{}                                                       // Pointer to a DTO struct. This struct should be kept in a closure along with the Run func (makes Run like a generic function)
	Secure bool                                                              // Request token will be authenticated iff this is true
	Run    func(*http.ResponseWriter, *http.Request, *config.User) *Response // Run the controller action.
}
