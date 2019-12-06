package api

import (
	"fmt"
	"net/http"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"

	"encoding/json"
)

// HandlerFactory converts web controllers into HandlerFuncs.
type HandlerFactory struct {
	authenticator *services.Authenticator
	config        *config.Config
}

// NewHandlerFactory constructs a new HandlerFactory.
func NewHandlerFactory(authenticator *services.Authenticator, config *config.Config) *HandlerFactory {
	return &HandlerFactory{
		authenticator: authenticator,
		config:        config,
	}
}

func (factory *HandlerFactory) NewHandler(controller *Controller) http.HandlerFunc {
	// Convert the controller into a binary controller
	b := func(w http.ResponseWriter, r *http.Request, u *config.User) *Response {

		context := controller.Make()

		d := json.NewDecoder(r.Body)
		err := d.Decode(context.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return nil
		}

		return context.Run(u, w, r)
	}
	return factory.NewBinaryHandler(&BinaryController{
		Run:    b,
		Secure: controller.Secure,
	})
}

func (factory *HandlerFactory) NewBinaryHandler(controller *BinaryController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", factory.config.AccessControlAllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Set-Cookie, DNT, Sec-Fetch-Mode")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Send no more to preflight requests
		if r.Method == "OPTIONS" {
			return
		}

		var response *Response
		var user *config.User

		// Authenticate!!!
		if controller.Secure {

			cookie, err := r.Cookie("token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user = factory.authenticator.AuthenticateFromJWT(cookie.Value)
			if user == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		response = controller.Run(w, r, user)
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
