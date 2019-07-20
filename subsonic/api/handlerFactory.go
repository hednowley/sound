package api

import (
	"io/ioutil"
	"net/http"
	"net/url"

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

func authenticate(params *url.Values, authenticator *services.Authenticator) bool {

	var auth bool
	var salt string
	var token string
	var password string

	username := params.Get("u")
	if len(username) == 0 {
		auth = false
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
		auth = authenticator.AuthenticateFromPassword(username, password)
		goto respond
	}

	token = params.Get("t")
	salt = params.Get("s")
	auth = authenticator.AuthenticateFromToken(username, salt, token)
	goto respond

respond:
	if !auth {
		log.Infof("Bad login attempt (%v).", username)
	}
	return auth

}

func (factory *HandlerFactory) PublishHandler(handler Handler) http.HandlerFunc {
	b := func(params url.Values, w *http.ResponseWriter, r *http.Request) *Response {
		return handler(params)
	}
	return factory.PublishBinaryHandler(b)
}

func (factory *HandlerFactory) PublishBinaryHandler(handler BinaryHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var response *Response

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			message := fmt.Sprintf("Error reading body: %v", err.Error())
			response = NewErrorReponse(dto.Generic, message)
		}

		//log.Debugf("\nREQUEST: %v\n%v\n", r.URL, string(body))

		params := parseParams(r.URL.Query(), body)

		format := parseResponseFormat(params.Get("f"))
		if format == nil {
			response = NewErrorReponse(dto.Generic, "Unknown format")
			format = &defaultFormat
			goto respond
		}

		if !authenticate(&params, factory.authenticator) {
			response = NewErrorReponse(dto.WrongCredentials, "Wrong username or password.")
			goto respond
		}

		response = handler(params, &w, r)
		if response != nil {
			goto respond
		}

		return

	respond:
		if *format == xmlFormat {
			s := serialiseToXML(response)
			fmt.Fprint(w, s)
			//log.Debugf("\nRESPONSE: %v\n%v\n", r.URL, s)
			return
		}

		if *format == jsonFormat {
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			s := serialiseToJSON(response)
			fmt.Fprint(w, s)
			//log.Debugf("\nRESPONSE: %v\n%v\n", r.URL, s)
			return
		}
	}
}
