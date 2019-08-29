package controller_test

import (
	"testing"

	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/controller"
	"github.com/hednowley/sound/api/dto"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

func TestGoodCredentials(t *testing.T) {
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

	cred, ok := c.Input.(*dto.Credentials)
	if !ok {
		t.Error()
	}

	cred.Username = "billy"
	cred.Password = "apple tart!!!"

	// Try and login with valid credentials
	r := c.Run(&config.User{})

	// Assert that the correct JWT was returned
	if r.Status != api.Success {
		t.Error()
	}

	token, ok := r.Body.(*dto.Token)
	if !ok {
		t.Error()
	}

	if token.Token != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1IjoiYmlsbHkifQ.YZyb0AaonWzRbrDRXc1sw4Y7BKYHtoR33NUmxP6iFSE" {
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

	cred, ok := c.Input.(*dto.Credentials)
	if !ok {
		t.Error()
	}

	cred.Username = ""
	cred.Password = ""

	// Try and login
	r := c.Run(&config.User{})

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

	cred, ok := c.Input.(*dto.Credentials)
	if !ok {
		t.Error()
	}

	cred.Username = "sdfsdfd"
	cred.Password = "ffff"

	// Try and login
	r := c.Run(&config.User{})

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
