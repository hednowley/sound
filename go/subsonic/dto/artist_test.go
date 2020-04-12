package dto

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
)

func TestArtist(t *testing.T) {

	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	GenerateAlbums(6, nil, artist, nil)

	DTO := NewArtist(artist)

	xml := `
	<artist id="1" name="artist1" coverArt="1.jpg" albumCount="6" duration="1"></artist>
	`

	json := `
	{
		"id":"1",
		"name":"artist1",
		"coverArt":"1.jpg",
		"albumCount":6,
		"duration":1
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestArtistWithoutArt(t *testing.T) {

	artist := GenerateArtist(1, nil)
	GenerateAlbums(6, nil, artist, nil)

	DTO := NewArtist(artist)

	xml := `
	<artist id="1" name="artist1" albumCount="6" duration="1"></artist>
	`

	json := `
	{
		"id":"1",
		"name":"artist1",
		"albumCount":6,
		"duration":1
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestArtistWithoutAlbums(t *testing.T) {

	art := GenerateArt(1)
	artist := GenerateArtist(1, art)

	DTO := NewArtist(artist)

	xml := `
	<artist id="1" name="artist1" coverArt="1.jpg" albumCount="0" duration="1"></artist>
	`

	json := `
	{
		"id":"1",
		"name":"artist1",
		"coverArt":"1.jpg",
		"albumCount":0,
		"duration":1
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestArtistWithAlbums(t *testing.T) {

	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	albums := GenerateAlbums(6, nil, artist, art)

	DTO := NewArtistWithAlbums(artist, albums)

	albumXML := ""
	for _, a := range albums {
		m, _ := xml.Marshal(NewAlbum(a))
		albumXML += string(m)
	}

	xml := fmt.Sprintf(`
	<artist id="1" name="artist1" coverArt="1.jpg" albumCount="6" duration="1">%v</artist>
	`, albumXML)

	albumJSON := make([]string, len(albums))
	for i, a := range albums {
		m, _ := json.Marshal(NewAlbum(a))
		albumJSON[i] = string(m)
	}

	json := fmt.Sprintf(`
	{
		"id":"1",
		"name":"artist1",
		"coverArt":"1.jpg",
		"albumCount":6,
		"album":[%v],
		"duration":1
	}
	`, strings.Join(albumJSON, ","))

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
