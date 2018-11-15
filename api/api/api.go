// Package api is stuff.
package api

import (
	"net/http"
)

// Controller is a web controller.
// It accepts a data-transfer object and returns an unserialised Response.
type Controller struct {
	Input  interface{}      // Pointer to a DTO struct. This struct should be kept in a closure along with the Run func (makes Run like a generic function)
	Secure bool             // Request token will be authenticated iff this is true
	Run    func() *Response // Run the controller action.
}

// BinaryController is a low-level web controller.
// It accepts a data-transfer object, a ResponseWriter and a Request.
// It returns a nil pointer to indicate that no further response is needed,
// otherwise it returns an unserialised Response.
type BinaryController struct {
	Input  interface{}                                         // Pointer to a DTO struct. This struct should be kept in a closure along with the Run func (makes Run like a generic function)
	Secure bool                                                // Request token will be authenticated iff this is true
	Run    func(*http.ResponseWriter, *http.Request) *Response // Run the controller action.
}
