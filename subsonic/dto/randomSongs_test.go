package dto

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
)

func TestRandomSongs(t *testing.T) {

	genre := GenerateGenre(1)
	art := GenerateArt(1)
	artist := GenerateArtist(1, art)
	album := GenerateAlbum(1, genre, artist, art)
	songs := GenerateSongs(5, genre, album, art)

	DTO := NewRandomSongs(songs)

	innerXML := ""
	for _, s := range songs {
		m, _ := xml.Marshal(NewSong(s))
		innerXML += string(m)
	}

	xml := fmt.Sprintf(`
	<randomSongs>%v</randomSongs>
	`, innerXML)

	innerJSON := make([]string, 5)
	for i, s := range songs {
		m, _ := json.Marshal(NewSong(s))
		innerJSON[i] = string(m)
	}

	json := fmt.Sprintf(`
	{
		"song":[%v]
	}
	`, strings.Join(innerJSON, ","))

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
