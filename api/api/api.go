package api

import (
	"net/http"
)

// Controller is a web controller.
// It accepts a data-transfer object and returns an unserialised Response.
type Controller struct {
	Input  interface{}      // Pointer to a DTO struct. This struct should be kept in a closure along with the Run field (makes Run like a generic function)
	Secure bool             // Request token will be authenticated iff this is true
	Run    func() *Response // Run the controller action.
}

// BinaryHandler is a low-level web controller action.
// It accepts a set of parameters, a ResponseWriter and a Request.
// It returns a nil pointer to indicate that no further response is needed,
// otherwise it returns an unserialised Response.
type BinaryHandler struct {
	Input  interface{}
	Secure bool
	Run    func(*http.ResponseWriter, *http.Request) *Response
}
