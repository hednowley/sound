package api

import (
	"fmt"
	"net/http"

	"github.com/hednowley/sound/services"

	"encoding/json"
)

// HandlerFactory converts friendly handlers into HandlerFuncs.
type HandlerFactory struct {
	authenticator *services.Authenticator
}

// NewHandlerFactory constructs a new HandlerFactory.
func NewHandlerFactory(authenticator *services.Authenticator) *HandlerFactory {
	return &HandlerFactory{
		authenticator: authenticator,
	}
}

func (factory *HandlerFactory) PublishHandler(handler *Controller) http.HandlerFunc {
	b := func(w *http.ResponseWriter, r *http.Request) *Response {
		return handler.Run()
	}
	return factory.PublishBinaryHandler(BinaryHandler{
		Input: handler.Input,
		Run:   b,
	})
}

func (factory *HandlerFactory) PublishBinaryHandler(handler BinaryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var response *Response

		d := json.NewDecoder(r.Body)
		err := d.Decode(&handler.Input)
		if err != nil {
			response = NewErrorReponse("Bad request.")
			goto respond
		}

		// Authenticate!!!
		if handler.Secure {

		}

		response = handler.Run(&w, r)
		if response != nil {
			goto respond
		}

		return

	respond:

		var body string
		data, err := json.Marshal(response.Body)
		if err != nil {
			body = fmt.Sprintf(`{ "status": "error", "message": "%v"}`, err)
		} else {
			var s string
			if response.Status == Success {
				s = "success"
			} else if response.Status == Fail {
				s = "fail"
			} else if response.Status == Error {
				s = "error"
			}

			body = fmt.Sprintf(`{ "status": "%v", "data": %v}`, s, string(data))
		}

		fmt.Fprint(w, body)
	}
}
