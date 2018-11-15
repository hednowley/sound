package api

import (
	"fmt"
	"net/http"

	"github.com/hednowley/sound/services"

	"encoding/json"
)

// HandlerFactory converts web controllers into HandlerFuncs.
type HandlerFactory struct {
	authenticator *services.Authenticator
}

// NewHandlerFactory constructs a new HandlerFactory.
func NewHandlerFactory(authenticator *services.Authenticator) *HandlerFactory {
	return &HandlerFactory{
		authenticator: authenticator,
	}
}

func (factory *HandlerFactory) NewHandler(controller *Controller) http.HandlerFunc {
	// Convert the controller into a binary controller
	b := func(w *http.ResponseWriter, r *http.Request) *Response {
		return controller.Run()
	}
	return factory.NewBinaryHandler(&BinaryController{
		Input: controller.Input,
		Run:   b,
	})
}

func (factory *HandlerFactory) NewBinaryHandler(controller *BinaryController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var response *Response

		d := json.NewDecoder(r.Body)
		err := d.Decode(&controller.Input)
		if err != nil {
			response = NewErrorReponse("Bad request.")
			goto respond
		}

		// Authenticate!!!
		if controller.Secure {

		}

		response = controller.Run(&w, r)
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
