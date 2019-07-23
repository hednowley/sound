package dto

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
)

func TestAlbum(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, art)
	GenerateSongs(11, genre, album, art)

	DTO := NewAlbum(album, false)

	xml := `
	<album id="1" name="album1" artist="artist1" artistId="1" coverArt="1.jpg" songCount="11" duration="1" created="2001-08-15T00:00:00Z" year="1901" genre="genre1"></album>
	`

	json := `
	{
		"id":"1",
		"name":"album1",
		"artist":"artist1",
		"artistId":"1",
		"coverArt":"1.jpg",
		"songCount":11,
		"duration":1,
		"created":"2001-08-15T00:00:00Z",
		"year":1901,
		"genre":"genre1"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlbumWithoutArt(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, nil)
	GenerateSongs(11, genre, album, art)

	DTO := NewAlbum(album, false)

	xml := `
	<album id="1" name="album1" artist="artist1" artistId="1" songCount="11" duration="1" created="2001-08-15T00:00:00Z" year="1901" genre="genre1"></album>
	`

	json := `
	{
		"id":"1",
		"name":"album1",
		"artist":"artist1",
		"artistId":"1",
		"songCount":11,
		"duration":1,
		"created":"2001-08-15T00:00:00Z",
		"year":1901,
		"genre":"genre1"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlbumWithoutGenre(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, nil, artist, art)
	GenerateSongs(11, genre, album, art)

	DTO := NewAlbum(album, false)

	xml := `
	<album id="1" name="album1" artist="artist1" artistId="1" coverArt="1.jpg" songCount="11" duration="1" created="2001-08-15T00:00:00Z" year="1901"></album>
	`

	json := `
	{
		"id":"1",
		"name":"album1",
		"artist":"artist1",
		"artistId":"1",
		"coverArt":"1.jpg",
		"songCount":11,
		"duration":1,
		"created":"2001-08-15T00:00:00Z",
		"year":1901
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlbumWithoutYear(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, art)
	GenerateSongs(11, genre, album, art)

	album.Year = 0

	DTO := NewAlbum(album, false)

	xml := `
	<album id="1" name="album1" artist="artist1" artistId="1" coverArt="1.jpg" songCount="11" duration="1" created="2001-08-15T00:00:00Z" genre="genre1"></album>
	`

	json := `
	{
		"id":"1",
		"name":"album1",
		"artist":"artist1",
		"artistId":"1",
		"coverArt":"1.jpg",
		"songCount":11,
		"duration":1,
		"created":"2001-08-15T00:00:00Z",
		"genre":"genre1"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlbumWithSongs(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, art)
	songs := GenerateSongs(11, genre, album, art)

	DTO := NewAlbum(album, true)

	songXML := ""
	for _, s := range songs {
		sd := NewSong(s)
		m, _ := xml.Marshal(sd)
		songXML += string(m)
	}

	xml := fmt.Sprintf(`
	<album id="1" name="album1" artist="artist1" artistId="1" coverArt="1.jpg" songCount="11" duration="1" created="2001-08-15T00:00:00Z" year="1901" genre="genre1">%v</album>
	`, songXML)

	songJSON := make([]string, len(songs))
	for i, s := range songs {
		sd := NewSong(s)
		m, _ := json.Marshal(sd)
		songJSON[i] = string(m)
	}

	json := fmt.Sprintf(`
	{
		"id":"1",
		"name":"album1",
		"artist":"artist1",
		"artistId":"1",
		"coverArt":"1.jpg",
		"songCount":11,
		"duration":1,
		"created":"2001-08-15T00:00:00Z",
		"year":1901,
		"genre":"genre1",
		"song":[%v]
	}
	`, strings.Join(songJSON, ","))

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
