package dto

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/hednowley/sound/database"
)

func TestArtist(t *testing.T) {

	artist, _ := database.NewMock().GetArtist(1)
	DTO := NewArtist(artist)

	xml := `
	<artist id="1" name="artist_1" coverArt="art_1.png" albumCount="3" duration="600"></artist>
	`

	json := `
	{"id":"1","name":"artist_1","coverArt":"art_1.png","albumCount":3,"duration":600}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestArtistWithoutArt(t *testing.T) {

	artist, _ := database.NewMock().GetArtist(6)
	DTO := NewArtist(artist)

	xml := `
	<artist id="6" name="artist_without_art" albumCount="1" duration="360"></artist>
	`

	json := `
	{"id":"6","name":"artist_without_art","albumCount":1,"duration":360}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestArtistWithoutAlbums(t *testing.T) {

	artist, _ := database.NewMock().GetArtist(7)
	DTO := NewArtist(artist)

	xml := `
	<artist id="7" name="artist_without_albums" albumCount="0" duration="0"></artist>
	`

	json := `
	{"id":"7","name":"artist_without_albums","albumCount":0,"duration":0}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestArtistWithAlbums(t *testing.T) {

	db := database.NewMock()

	artist, err := db.GetArtist(1)
	if err != nil {
		t.Error(err)
		return
	}

	albums, err := db.GetAlbumsByArtist(1)
	if err != nil {
		t.Error(err)
		return
	}

	DTO := NewArtistWithAlbums(artist, albums)

	albumXML := ""
	for _, a := range albums {
		m, _ := xml.Marshal(NewAlbum(&a))
		albumXML += string(m)
	}

	xml := fmt.Sprintf(`
	<artist id="1" name="artist_1" coverArt="art_1.png" albumCount="3" duration="600">%v</artist>
	`, albumXML)

	albumJSON := make([]string, len(albums))
	for i, a := range albums {
		m, _ := json.Marshal(NewAlbum(&a))
		albumJSON[i] = string(m)
	}

	json := fmt.Sprintf(`
	{
		"id":"1",
		"name":"artist_1",
		"coverArt":"art_1.png",
		"albumCount":3,
		"album":[%v],
		"duration":600
	}
	`, strings.Join(albumJSON, ","))

	err = CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
