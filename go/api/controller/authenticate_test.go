package controller_test

import (
	"net/http/httptest"
	"testing"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/controller"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

func TestGoodCredentials(t *testing.T) {

	rr := httptest.NewRecorder()

	cfg := config.Config{
		Secret: "secret",
		Users: []config.User{
			{
				Username: "billy",
				Password: "apple tart!!!",
				Email:    "billy@bigbugs.com",
			},
		},
	}
	a := services.NewAuthenticator(&cfg)
	c := controller.NewAuthenticateController(a)

	context := c.Make()

	cred, ok := context.Body.(*dto.Credentials)
	if !ok {
		t.Error()
	}

	cred.Username = "billy"
	cred.Password = "apple tart!!!"

	// Try and login with valid credentials
	r := context.Run(&config.User{}, rr, nil)

	// Assert that the correct JWT was returned
	if r.Status != api.Success {
		t.Error()
	}

}

func TestEmptyCredentials(t *testing.T) {
	cfg := config.Config{
		Secret: "secret",
		Users: []config.User{
			{
				Username: "billy",
				Password: "apple tart!!!",
				Email:    "billy@bigbugs.com",
			},
		},
	}
	a := services.NewAuthenticator(&cfg)
	c := controller.NewAuthenticateController(a)

	context := c.Make()

	context.Body = dto.Credentials{}

	// Try and login
	r := context.Run(&config.User{}, nil, nil)

	if r.Status != api.Error {
		t.Error()
	}

	message, ok := r.Body.(string)
	if !ok {
		t.Error()
	}

	if message != "Bad credentials." {
		t.Error()
	}
}

func TestBadCredentials(t *testing.T) {
	cfg := config.Config{
		Secret: "secret",
		Users: []config.User{
			{
				Username: "billy",
				Password: "apple tart!!!",
				Email:    "billy@bigbugs.com",
			},
		},
	}
	a := services.NewAuthenticator(&cfg)
	c := controller.NewAuthenticateController(a)

	context := c.Make()

	context.Body = dto.Credentials{
		Username: "sdfsdfd",
		Password: "ffff",
	}

	// Try and login
	r := context.Run(&config.User{}, nil, nil)

	if r.Status != api.Error {
		t.Error()
	}

	message, ok := r.Body.(string)
	if !ok {
		t.Error()
	}

	if message != "Bad credentials." {
		t.Error()
	}
}
