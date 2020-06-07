package api

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
	"github.com/hednowley/sound/subsonic/dto"

	"encoding/hex"
	"fmt"
	"strings"

	log "github.com/cihub/seelog"
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

func authenticate(params *url.Values, authenticator *services.Authenticator) *config.User {

	var user *config.User
	var salt string
	var token string
	var password string

	username := params.Get("u")
	if len(username) == 0 {
		goto respond
	}

	password = params.Get("p")
	if len(password) != 0 {
		if strings.HasPrefix(password, "enc:") {
			encoded := strings.TrimPrefix(password, "enc:")
			bytes, err := hex.DecodeString(encoded)
			if err == nil {
				password = string(bytes)
			}
		}
		user = authenticator.AuthenticateFromPassword(username, password)
		goto respond
	}

	token = params.Get("t")
	salt = params.Get("s")
	user = authenticator.AuthenticateFromToken(username, salt, token)
	goto respond

respond:
	if user == nil {
		log.Infof("Bad login attempt (%v).", username)
	}
	return user

}

// PublishHandler converts a friendly handler into a HandlerFunc.
func (factory *HandlerFactory) PublishHandler(handler Handler) http.HandlerFunc {
	b := func(params url.Values, w *http.ResponseWriter, r *http.Request, c *HandlerContext) *Response {
		return handler(params, c)
	}
	return factory.PublishBinaryHandler(b)
}

// PublishBinaryHandler converts a friendly binary handler into a HandlerFunc.
func (factory *HandlerFactory) PublishBinaryHandler(handler BinaryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var response *Response
		var format *responseFormat
		var params url.Values
		var user *config.User
		var context HandlerContext

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			message := fmt.Sprintf("Error reading body: %v", err.Error())
			response = NewErrorReponse(dto.Generic, message)
			goto respond
		}

		params = parseParams(r.URL.Query(), body)

		format = parseResponseFormat(params.Get("f"))
		if format == nil {
			response = NewErrorReponse(dto.Generic, "Unknown format")
			format = &defaultFormat
			goto respond
		}

		// All endpoints require authentication
		user = authenticate(&params, factory.authenticator)
		if user == nil {
			response = NewErrorReponse(dto.WrongCredentials, "Wrong username or password.")
			goto respond
		}

		context.User = user
		response = handler(params, &w, r, &context)
		if response != nil {
			goto respond
		}

		return

	respond:
		if *format == xmlFormat {
			s := serialiseToXML(response)
			fmt.Fprint(w, s)
			return
		}

		if *format == jsonFormat {
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			s := serialiseToJSON(response)
			fmt.Fprint(w, s)
			return
		}
	}
}
