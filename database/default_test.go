package database_test

import (
	"testing"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/database"
)

// Putters
func TestPutArtist(t *testing.T) {
	m := database.NewMock()

	// Update existing record
	m.PutArtist(&dao.Artist{
		ID:   2,
		Name: "artist_2X",
	})

	a := m.GetArtist(2, false, false)
	if a.Name != "artist_2X" {
		t.Error("Not modified")
	}

	// Add new record with ID
	m.PutArtist(&dao.Artist{
		ID:   666,
		Name: "heggarty",
	})
	a = m.GetArtist(666, false, false)
	if a.Name != "heggarty" {
		t.Error("Not modified")
	}

	// Add new record without ID
	a = &dao.Artist{
		Name: "sdff",
	}
	m.PutArtist(a)

	// Should have mutated with a new ID
	if a.ID != 10001 {
		t.Error("Wrong ID")
	}

	a = m.GetArtist(10001, false, false)
	if a.Name != "sdff" {
		t.Error("Not modified")
	}
}

/*
func TestPutArtistByName(t *testing.T) {
	m := database.NewMock()

	// Get existing artist
	a := m.PutArtistByName("artist_2")
	if a.ID != 2 {
		t.Error("Wrong ID")
	}

	// Get new artist
	a = m.PutArtistByName("heggarty")
	if a.ID != 10001 {
		t.Error("Wrong ID")
	}

	a = m.GetArtist(10001, false, false)
	if a.Name != "heggarty" {
		t.Error("Not added")
	}
}
*/

func TestPutAlbumByNameAndArtist(t *testing.T) {
	m := database.NewMock()

	// Put existing album
	a := m.PutAlbumByNameAndArtist("album_2", "artist_1")
	if a.ID != 2 {
		t.Error()
	}

	// Put album with existing name but new artist
	a = m.PutAlbumByNameAndArtist("album_2", "sdffsd")
	if a.ID != 10001 || a.ArtistID != 10001 {
		t.Error()
	}

	a2 := m.GetArtist(10001, true, false)
	if a2.Name != "sdffsd" || a2.Albums[0].ID != 10001 {
		t.Error()
	}

	// Put album with new name but existing artist
	a = m.PutAlbumByNameAndArtist("ghghddh", "artist_2")
	if a.ID != 10002 || a.ArtistID != 2 {
		t.Error()
	}

	// Put new album and artist
	a = m.PutAlbumByNameAndArtist("sdfgyifs", "sduiru")
	if a.ID != 10003 || a.ArtistID != 10002 {
		t.Error()
	}

	a2 = m.GetArtist(10002, true, false)
	if a2.Name != "sduiru" || a2.Albums[0].ID != 10003 {
		t.Error()
	}
}

func TestPutGenreByName(t *testing.T) {
	m := database.NewMock()

	// Get existing genre
	g := m.PutGenreByName("genre_2")
	if g.ID != 2 || g.Name != "genre_2" {
		t.Error("Wrong ID")
	}

	// Get new genre
	g = m.PutGenreByName("crustpunk")
	if g.ID != 10001 || g.Name != "crustpunk" {
		t.Error("Wrong ID")
	}
}

func TestPutSong(t *testing.T) {
	m := database.NewMock()

	// Update existing record
	m.PutSong(&dao.Song{
		ID:      2,
		AlbumID: 2,          // Different
		Path:    "new_path", // Different
		Title:   "djndjnd",  // Different
	})

	s := m.GetSong(2, false, false, false, false)
	if s.AlbumID != 2 || s.Path != "new_path" || s.Title != "djndjnd" {
		t.Error()
	}

	// Add new record
	s = &dao.Song{
		AlbumID: 2,
		Path:    "new_path",
		Title:   "djndjnd",
	}
	m.PutSong(s)

	if s.ID != 10001 {
		t.Error()
	}

	s = m.GetSong(10001, false, false, false, false)
	if s.AlbumID != 2 || s.Path != "new_path" || s.Title != "djndjnd" {
		t.Error()
	}
}

func TestDelete(t *testing.T) {
	m := database.NewMock()
	m.Empty()

	if len(m.GetArtists()) > 0 {
		t.Error()
	}
}
