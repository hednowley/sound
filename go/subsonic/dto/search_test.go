package dto

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/hednowley/sound/dao"
)

func TestSearch2(t *testing.T) {

	genres := GenerateGenres(3)
	arts := GenerateArts(3)

	artist1 := GenerateArtist(1, arts[0])
	artist2 := GenerateArtist(2, arts[1])
	artists := []*dao.Artist{artist1, artist2}

	album1 := GenerateAlbum(1, genres[0], artist1, arts[0])
	album2 := GenerateAlbum(2, genres[0], artist2, arts[0])
	album3 := GenerateAlbum(3, genres[0], artist2, arts[0])
	albums := []*dao.Album{album1, album2, album3}

	song1 := GenerateSong(1, genres[0], album1, arts[0])
	song2 := GenerateSong(5, genres[1], album2, arts[1])
	song3 := GenerateSong(10, genres[2], album3, arts[2])
	songs := []*dao.Song{song1, song2, song3}

	DTO := NewSearch2Response(artists, albums, songs)

	artistsJSON := make([]string, 0)
	albumsJSON := make([]string, 0)
	songsJSON := make([]string, 0)
	innerXML := ""

	for _, a := range artists {
		m, _ := xml.Marshal(NewArtist(a))
		innerXML += string(m)

		m, _ = json.Marshal(NewArtist(a))
		artistsJSON = append(artistsJSON, string(m))
	}

	for _, a := range albums {
		m, _ := xml.Marshal(NewAlbum(a))
		innerXML += string(m)

		m, _ = json.Marshal(NewAlbum(a))
		albumsJSON = append(albumsJSON, string(m))
	}

	for _, s := range songs {
		m, _ := xml.Marshal(NewSong(s))
		innerXML += string(m)

		m, _ = json.Marshal(NewSong(s))
		songsJSON = append(songsJSON, string(m))
	}

	xml := fmt.Sprintf(`
	<searchResult2>%s</searchResult2>
	`, innerXML)

	json := fmt.Sprintf(`
	{
		"artist":[%s],
		"album":[%s],
		"song":[%s]
	}
	 `,
		strings.Join(artistsJSON, ","),
		strings.Join(albumsJSON, ","),
		strings.Join(songsJSON, ","))

	err := CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
