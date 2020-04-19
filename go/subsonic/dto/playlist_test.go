package dto

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/hednowley/sound/database"
)

func TestPlaylist(t *testing.T) {

	db := database.NewMock()

	playlist, err := db.GetPlaylist(1)
	if err != nil {
		t.Error(err)
		return
	}

	songs, err := db.GetPlaylistSongs(1)
	if err != nil {
		t.Error(err)
		return
	}

	DTO := NewPlaylist(playlist, songs)

	entryXML := ""
	for _, s := range songs {
		m, _ := xml.Marshal(newPlaylistEntry(&s))
		entryXML += string(m)
	}

	xml := fmt.Sprintf(`
	<playlist id="1" name="playlist_1" comment="comment_1" owner="ned" public="true" songCount="4" duration="480" created="2018-06-12T11:11:11+01:00" changed="2018-06-12T11:11:11+01:00">%v</playlist>
	`, entryXML)

	entryJSON := []string{}
	for _, s := range songs {
		m, _ := json.Marshal(newPlaylistEntry(&s))
		entryJSON = append(entryJSON, string(m))
	}

	json := fmt.Sprintf(`
	{
		"id":"1",
		"name":"playlist_1",
		"comment":"comment_1",
		"owner":"ned",
		"public":true,
		"songCount":4,
		"duration":480,
		"created":"2018-06-12T11:11:11+01:00",
		"changed":"2018-06-12T11:11:11+01:00",
		"entry":[%v]
	}
	 `, strings.Join(entryJSON, ","))

	err = CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestPlaylistWithoutComment(t *testing.T) {

	db := database.NewMock()

	playlist, err := db.GetPlaylist(1)
	if err != nil {
		t.Error(err)
		return
	}

	playlist.Comment = ""

	songs, err := db.GetPlaylistSongs(1)
	if err != nil {
		t.Error(err)
		return
	}

	DTO := NewPlaylist(playlist, songs)

	entryXML := ""
	for _, s := range songs {
		m, _ := xml.Marshal(newPlaylistEntry(&s))
		entryXML += string(m)
	}

	xml := fmt.Sprintf(`
	<playlist id="1" name="playlist_1" comment="" owner="ned" public="true" songCount="4" duration="480" created="2018-06-12T11:11:11+01:00" changed="2018-06-12T11:11:11+01:00">%v</playlist>
	`, entryXML)

	entryJSON := []string{}
	for _, s := range songs {
		m, _ := json.Marshal(newPlaylistEntry(&s))
		entryJSON = append(entryJSON, string(m))
	}

	json := fmt.Sprintf(`
	{
		"id":"1",
		"name":"playlist_1",
		"comment":"",
		"owner":"ned",
		"public":true,
		"songCount":4,
		"duration":480,
		"created":"2018-06-12T11:11:11+01:00",
		"changed":"2018-06-12T11:11:11+01:00",
		"entry":[%v]
	}
	 `, strings.Join(entryJSON, ","))

	err = CheckSerialisation(DTO, xml, json)
	if err != nil {
		t.Error(err.Error())
	}
}
