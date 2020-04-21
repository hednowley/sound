package database_test

import (
	"testing"

	"github.com/hednowley/sound/database"
)

// Putters
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

func TestPutAlbumByAttributest(t *testing.T) {
	m := database.NewMock()

	// Put existing album
	albumID, err := m.PutAlbumByAttributes("album_without_art", "artist_without_art", "")
	if err != nil {
		t.Error(err)
		return
	}
	if albumID != 2 {
		t.Error()
		return
	}

	// Put album with existing name but new artist
	albumID, err = m.PutAlbumByAttributes("album_without_art", "sdffsd", "")
	if err != nil {
		t.Error(err)
		return
	}
	if albumID != 10001 {
		t.Error()
		return
	}

	a, err := m.GetAlbum(albumID)
	if err != nil {
		t.Error(err)
		return
	}
	if a == nil || a.ArtistID != 10001 {
		t.Error()
		return
	}

	// Put album with new name but existing artist
	albumID, err = m.PutAlbumByAttributes("ghghddh", "artist_2", "")
	if albumID != 10002 {
		t.Error()
		return
	}
	a, err = m.GetAlbum(albumID)
	if err != nil {
		t.Error(err)
		return
	}
	if a == nil || a.ArtistID != 2 {
		t.Error()
		return
	}

	// Put new album and artist
	albumID, err = m.PutAlbumByAttributes("sdfgyifs", "sduiru", "")
	if albumID != 10003 {
		t.Error()
		return
	}
	a, err = m.GetAlbum(albumID)
	if err != nil {
		t.Error(err)
		return
	}
	if a == nil || a.ArtistID != 10002 {
		t.Error()
		return
	}

	// Put existing album with new disambiguator
	albumID, err = m.PutAlbumByAttributes("album_2", "artist_1", "d1")
	if albumID != 10004 {
		t.Error()
		return
	}

	// Put same album again
	albumID, err = m.PutAlbumByAttributes("album_2", "artist_1", "d1")
	if albumID != 10004 {
		t.Error()
		return
	}

	// Put same album again with another new disambiguator
	albumID, err = m.PutAlbumByAttributes("album_2", "artist_1", "d2")
	if albumID != 10005 {
		t.Error()
		return
	}

	// Put album with same artist but different capitalisation
	albumID, err = m.PutAlbumByAttributes("album_3", "ArtisT_1", "")
	if albumID != 10006 {
		t.Error()
		return
	}
	a, err = m.GetAlbum(albumID)
	if err != nil {
		t.Error(err)
		return
	}
	if a == nil || a.ArtistID != 1 {
		t.Error()
		return
	}
}

func TestPutGenreByName(t *testing.T) {
	m := database.NewMock()

	// Get existing genre
	g, err := m.PutGenreByName("genre_2")
	if err != nil || g != 2 {
		t.Error("Wrong ID")
	}

	// Get new genre
	g, err = m.PutGenreByName("crustpunk")
	if err != nil || g != 10001 {
		t.Error("Wrong ID")
	}
}
