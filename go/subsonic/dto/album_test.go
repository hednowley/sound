package dto

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/hednowley/sound/database"
)

func TestAlbum(t *testing.T) {

	db := database.NewMock()
	conn, _ := db.GetConn()
	defer conn.Release()

	album, _ := db.GetAlbum(conn, 1)
	DTO := NewAlbum(album)

	xml := `
	<album id="1" name="album_1" artist="artist_1" artistId="1" coverArt="art_1.png" songCount="3" duration="360" created="2018-06-12T11:11:11Z" year="2000" genre="genre_1"></album>
	`

	json := `
	{
		"id":"1",
		"name":"album_1",
		"artist":"artist_1",
		"artistId":"1",
		"coverArt":"art_1.png",
		"songCount":3,
		"duration":360,
		"created":"2018-06-12T11:11:11Z",
		"year":2000,
		"genre":"genre_1"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlbumWithoutArt(t *testing.T) {

	db := database.NewMock()
	conn, _ := db.GetConn()
	defer conn.Release()

	album, _ := db.GetAlbum(conn, 2)
	DTO := NewAlbum(album)

	xml := `
	<album id="2" name="album_without_art" artist="artist_without_art" artistId="6" songCount="3" duration="360" created="2018-06-12T11:11:11Z" year="1997" genre="genre_1"></album>
	`

	json := `
	{
		"id":"2",
		"name":"album_without_art",
		"artist":"artist_without_art",
		"artistId":"6",
		"songCount":3,
		"duration":360,
		"created":"2018-06-12T11:11:11Z",
		"year":1997,
		"genre":"genre_1"
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlbumWithoutGenre(t *testing.T) {

	db := database.NewMock()
	conn, _ := db.GetConn()
	defer conn.Release()

	album, _ := db.GetAlbum(conn, 3)
	DTO := NewAlbum(album)

	xml := `
	<album id="3" name="album_without_genre" artist="artist_1" artistId="1" songCount="1" duration="120" created="2018-06-12T11:11:11Z" year="1964"></album>
	`

	json := `
	{
		"id":"3",
		"name":"album_without_genre",
		"artist":"artist_1",
		"artistId":"1",
		"songCount":1,
		"duration":120,
		"created":"2018-06-12T11:11:11Z",
		"year":1964
	}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlbumWithoutYear(t *testing.T) {

	db := database.NewMock()
	conn, _ := db.GetConn()
	defer conn.Release()

	album, _ := db.GetAlbum(conn, 4)
	DTO := NewAlbum(album)

	xml := `
	<album id="4" name="album_without_year" artist="artist_1" artistId="1" songCount="1" duration="120" created="2018-06-12T11:11:11Z" genre="genre_1"></album>
	`

	json := `
	{"id":"4","name":"album_without_year","artist":"artist_1","artistId":"1","songCount":1,"duration":120,"created":"2018-06-12T11:11:11Z","genre":"genre_1"}
	`

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlbumWithSongs(t *testing.T) {

	db := database.NewMock()
	conn, _ := db.GetConn()
	defer conn.Release()

	album, _ := db.GetAlbum(conn, 1)
	songs, _ := db.GetAlbumSongs(conn, 1)

	DTO := NewAlbumWithSongs(album, songs)

	songXML := ""
	for _, s := range songs {
		sd := NewSong(&s)
		m, _ := xml.Marshal(sd)
		songXML += string(m)
	}

	xml := fmt.Sprintf(`
	<album id="1" name="album_1" artist="artist_1" artistId="1" coverArt="art_1.png" songCount="3" duration="360" created="2018-06-12T11:11:11Z" year="2000" genre="genre_1">%v</album>
	`, songXML)

	songJSON := make([]string, len(songs))
	for i, s := range songs {
		sd := NewSong(&s)
		m, _ := json.Marshal(sd)
		songJSON[i] = string(m)
	}

	json := fmt.Sprintf(`
	{
		"id":"1",
		"name":"album_1",
		"artist":"artist_1",
		"artistId":"1",
		"coverArt":"art_1.png",
		"songCount":3,
		"duration":360,
		"created":"2018-06-12T11:11:11Z",
		"year":2000,
		"genre":"genre_1",
		"song":[%v]
	}
	`, strings.Join(songJSON, ","))

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
