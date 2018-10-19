package dto

import (
	"testing"

	"github.com/hednowley/sound/config"
)

func TestUser(t *testing.T) {

	user := config.User{
		Username: "dsauid",
		Email:    "sdfhjfsd@dsid.com",
	}

	DTO := NewUser(user)

	xml := `
	<user username="dsauid" email="sdfhjfsd@dsid.com" scrobblingEnabled="false" adminRole="true" settingsRole="true" downloadRole="true" uploadRole="true" playlistRole="true" coverArtRole="true" commentRole="true" podcastRole="true" streamRole="true" jukeboxRole="true" shareRole="true"></user>
	`

	json := `
	{
		"username":"dsauid",
		"email":"sdfhjfsd@dsid.com",
		"scrobblingEnabled":false,
		"adminRole":true,
		"settingsRole":true,
		"downloadRole":true,
		"uploadRole":true,
		"playlistRole":true,
		"coverArtRole":true,
		"commentRole":true,
		"podcastRole":true,
		"streamRole":true,
		"jukeboxRole":true,
		"shareRole":true
	 }
	 `

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
