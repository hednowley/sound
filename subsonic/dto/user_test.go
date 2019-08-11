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
	<user username="dsauid" email="sdfhjfsd@dsid.com" scrobblingEnabled="false" adminRole="true" settingsRole="true" downloadRole="true" uploadRole="true" playlistRole="true" coverArtRole="false" commentRole="false" podcastRole="false" streamRole="true" jukeboxRole="false" shareRole="false" videoConversionRole="false" folder="0"></user>
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
		"coverArtRole":false,
		"commentRole":false,
		"podcastRole":false,
		"streamRole":true,
		"jukeboxRole":false,
		"shareRole":false,
		"videoConversionRole":false,
		"folder":[0]
	 }
	 `

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
