package services_test

import (
	"testing"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/services"
)

var c = config.Config{
	Users: []config.User{
		config.User{
			Username: "billy",
			Password: "apple tart!!!",
			Email:    "billy@bigbugs.com",
		},
		config.User{
			Username: "tom tom",
			Password: "sfjksdfjk",
			Email:    "tom@bigbugs.com",
		},
	},
}

var a = services.NewAuthenticator(&c)

func TestPasswordAuth(t *testing.T) {

	// Should work
	if !a.AuthenticateFromPassword("billy", "apple tart!!!") {
		t.Error()
	}

	// Shouldn't work
	if a.AuthenticateFromPassword("billy2", "apple tart!!!") {
		t.Error()
	}

	if a.AuthenticateFromPassword("billy", "appletart!!!") {
		t.Error()
	}

	if a.AuthenticateFromPassword("", "") {
		t.Error()
	}
}

func TestTokenAuth(t *testing.T) {

	// Should work
	if !a.AuthenticateFromToken("tom tom", "saltysalt", "bbf729f2464585f1212e03519659b30a") {
		t.Error()
	}

	// Shouldn't work
	if a.AuthenticateFromToken("tom tom", "saltysalt", "bbf729f2464585f12e03519659b30a") {
		t.Error()
	}

	if a.AuthenticateFromToken("tom tom", "rocksaltt", "bbf729f2464585f1212e03519659b30a") {
		t.Error()
	}

	if a.AuthenticateFromToken("", "", "") {
		t.Error()
	}
}
