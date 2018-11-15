package api

import (
	"net/http"
)

// Handler is a web controller action.
// It accepts a set of parameters and returns an unserialised Response.
type Handler struct {
	Input  interface{}
	Worker func() *Response
}

func (h *Handler) Run() {
	//h.worker
}

// BinaryHandler is a low-level web controller action.
// It accepts a set of parameters, a ResponseWriter and a Request.
// It returns a nil pointer to indicate that no further response is needed,
// otherwise it returns an unserialised Response.
type BinaryHandler struct {
	Input  interface{}
	Worker func(*http.ResponseWriter, *http.Request) *Response
}
