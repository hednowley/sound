package provider_test

import (
	"testing"

	"github.com/hednowley/sound/ws"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/provider"
)

func TestAddOnlyScan(t *testing.T) {

	f := []*entities.FileInfo{
		// Existing song with different title
		&entities.FileInfo{
			Album:       "album_2",
			AlbumArtist: "artist_2",
			Path:        "path_2",
			Title:       "Y.M.C.A.",
		},
		// New song
		&entities.FileInfo{
			Album:       "new_album",
			AlbumArtist: "new_artist",
			Path:        "new_path",
			Title:       "new_title",
		},
	}

	p := provider.NewMockProvider("mock", f)
	p.SetScanID("scan1")
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{}, m)
	hub := ws.NewMockHub()
	scanner := provider.NewScanner([]provider.Provider{p}, dal, hub)

	// Scan without updating or deleting
	scanner.StartAllScans(false, false)

	// Only scan ID should change
	s := m.GetSong(2, false, false, false, false)
	if s.Title != "title_2" || s.ScanID != "scan1" {
		t.Error()
	}

	// Scan ID shouldn't change as this song wasn't provided
	s = m.GetSong(1, false, false, false, false)
	if s.ScanID != "" {
		t.Error()
	}

	// Scan ID shouldn't change as this song has another provider
	s = m.GetSong(22, false, false, false, false)
	if s.ScanID != "" {
		t.Error()
	}

	// New song should have been added
	s = m.GetSong(10001, false, false, false, false)
	if s == nil || s.ScanID != "scan1" || s.Title != "new_title" {
		t.Error()
	}

	// New artist should have been added
	artist := m.GetArtist(10001, false, false)
	if artist == nil || artist.Name != "new_artist" {
		t.Error()
	}

	// New album should have been added
	album := m.GetAlbum(10001, false, false, false)
	if album == nil || album.Name != "new_album" || album.ArtistID != 10001 {
		t.Error()
	}
}

func TestUpdateScan(t *testing.T) {

	f := []*entities.FileInfo{
		// Existing song with different info
		&entities.FileInfo{
			Album:       "album_2",
			AlbumArtist: "artist_1",
			Path:        "path_2",
			Title:       "Y.M.C.A.",  // Different
			Genre:       "Neurofunk", // Different
		},
		// New song
		&entities.FileInfo{
			Album:       "new_album",
			AlbumArtist: "new_artist",
			Path:        "new_path",
			Title:       "new_title",
		},
	}

	p := provider.NewMockProvider("mock", f)
	p.SetScanID("scan1")
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{}, m)
	hub := ws.NewMockHub()
	scanner := provider.NewScanner([]provider.Provider{p}, dal, hub)

	// Scan without updating or deleting
	scanner.StartAllScans(true, false)

	// Should change
	s := m.GetSong(2, false, true, false, false)
	if s.Title != "Y.M.C.A." || s.ScanID != "scan1" {
		t.Error()
	}

	// New genre should have been added
	genre := m.GetGenre("Neurofunk")
	if genre == nil || s.GenreID != genre.ID || s.Album.GenreID != genre.ID {
		t.Error()
	}

	// Scan ID shouldn't change as this song wasn't provided
	s = m.GetSong(1, false, false, false, false)
	if s.ScanID != "" {
		t.Error()
	}

	// Scan ID shouldn't change as this song has another provider
	s = m.GetSong(22, false, false, false, false)
	if s.ScanID != "" {
		t.Error()
	}

	// New song should have been added
	s = m.GetSong(10001, true, true, true, true)
	if s == nil || s.ScanID != "scan1" || s.Title != "new_title" {
		t.Error()
	}

	// New artist should have been added
	artist := m.GetArtist(10001, false, false)
	if artist == nil || artist.Name != "new_artist" {
		t.Error()
	}

	// New album should have been added
	album := m.GetAlbum(10001, false, false, false)
	if album == nil || album.Name != "new_album" || album.ArtistID != 10001 {
		t.Error()
	}

}
