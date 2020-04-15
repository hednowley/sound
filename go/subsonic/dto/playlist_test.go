package dto

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/hednowley/sound/dao"
)

func TestPlaylist(t *testing.T) {

	genres := GenerateGenres(3)
	arts := GenerateArts(3)

	artist1 := GenerateArtist(1, &arts[0])
	artist2 := GenerateArtist(2, &arts[1])

	album1 := GenerateAlbum(1, &genres[0], artist1, &arts[0])
	album2 := GenerateAlbum(2, &genres[0], artist2, &arts[0])
	album3 := GenerateAlbum(3, &genres[0], artist2, &arts[0])

	song1 := GenerateSong(1, &genres[0], album1, &arts[0])
	song2 := GenerateSong(5, &genres[1], album2, &arts[1])
	song3 := GenerateSong(10, &genres[2], album3, &arts[2])

	playlist := GeneratePlaylist(1, true)

	DTO := NewPlaylist(playlist, []dao.Song{*song1, *song3, *song2})

	entryXML := ""
	for _, a := range []*dao.Song{song1, song3, song2} {
		m, _ := xml.Marshal(newPlaylistEntry(a))
		entryXML += string(m)
	}

	xml := fmt.Sprintf(`
	<playlist id="1" name="playlist1" comment="comment1" owner="ned" public="true" songCount="3" duration="0" created="2002-01-15T05:01:00.000000001Z" changed="2001-08-10T00:30:00.000000009Z">%v</playlist>
	`, entryXML)

	entryJSON := make([]string, 3)
	for i, a := range []*dao.Song{song1, song3, song2} {
		m, _ := json.Marshal(newPlaylistEntry(a))
		entryJSON[i] = string(m)
	}

	json := fmt.Sprintf(`
	{
		"id":"1",
		"name":"playlist1",
		"comment":"comment1",
		"owner":"ned",
		"public":true,
		"songCount":3,
		"duration":0,
		"created":"2002-01-15T05:01:00.000000001Z",
		"changed":"2001-08-10T00:30:00.000000009Z",
		"entry":[%v]
	}
	 `, strings.Join(entryJSON, ","))

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestPlaylistWithoutComment(t *testing.T) {

	genres := GenerateGenres(3)
	arts := GenerateArts(3)

	artist1 := GenerateArtist(1, &arts[0])
	artist2 := GenerateArtist(2, &arts[1])

	album1 := GenerateAlbum(1, &genres[0], artist1, &arts[0])
	album2 := GenerateAlbum(2, &genres[0], artist2, &arts[0])
	album3 := GenerateAlbum(3, &genres[0], artist2, &arts[0])

	song1 := GenerateSong(1, &genres[0], album1, &arts[0])
	song2 := GenerateSong(5, &genres[1], album2, &arts[1])
	song3 := GenerateSong(10, &genres[2], album3, &arts[2])

	playlist := GeneratePlaylist(1, true)
	playlist.Comment = ""

	DTO := NewPlaylist(playlist, []dao.Song{*song1, *song3, *song2})

	entryXML := ""
	for _, a := range []*dao.Song{song1, song3, song2} {
		m, _ := xml.Marshal(newPlaylistEntry(a))
		entryXML += string(m)
	}

	xml := fmt.Sprintf(`
	<playlist id="1" name="playlist1" comment="" owner="ned" public="true" songCount="3" duration="0" created="2002-01-15T05:01:00.000000001Z" changed="2001-08-10T00:30:00.000000009Z">%v</playlist>
	`, entryXML)

	entryJSON := make([]string, 3)
	for i, a := range []*dao.Song{song1, song3, song2} {
		m, _ := json.Marshal(newPlaylistEntry(a))
		entryJSON[i] = string(m)
	}

	json := fmt.Sprintf(`
	{
		"id":"1",
		"name":"playlist1",
		"comment":"",
		"owner":"ned",
		"public":true,
		"songCount":3,
		"duration":0,
		"created":"2002-01-15T05:01:00.000000001Z",
		"changed":"2001-08-10T00:30:00.000000009Z",
		"entry":[%v]
	}
	 `, strings.Join(entryJSON, ","))

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
