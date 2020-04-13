package provider_test

import (
	"testing"

	"github.com/hednowley/sound/socket"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/entities"
	"github.com/hednowley/sound/provider"
)

func TestAddOnlyScan(t *testing.T) {

	f := []*entities.FileInfo{
		// Existing song with different title
		{
			Album:       "album_2",
			AlbumArtist: "artist_2",
			Path:        "path_2",
			Title:       "Y.M.C.A.",
		},
		// New song
		{
			Album:       "new_album",
			AlbumArtist: "new_artist",
			Path:        "new_path",
			Title:       "new_title",
		},
	}

	p := provider.NewMockProvider("mock", f)
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{}, m)
	hub := socket.NewMockHub()
	scanner := provider.NewScanner([]provider.Provider{p}, dal, hub)

	// Scan without updating or deleting
	scanner.StartAllScans(false, false)

	// Title shouldn't change
	s, err := dal.GetSong(2)
	if err != nil || s.Title != "title_2" {
		t.Error()
	}

	// Song shouldn't change as this song wasn't provided
	s, err = dal.GetSong(1)
	if err != nil {
		t.Error()
	}

	// Song shouldn't change as this song has another provider
	s, err = dal.GetSong(22)
	if err != nil {
		t.Error()
	}

	// New song should have been added
	s, err = dal.GetSong(10001)
	if err != nil || s == nil || s.Title != "new_title" {
		t.Error()
	}

	// New artist should have been added
	artist, err := dal.GetArtist(10001)
	if err != nil || artist == nil || artist.Name != "new_artist" {
		t.Error()
	}

	// New album should have been added
	album, err := dal.GetAlbum(10001)
	if err != nil || album == nil || album.Name != "new_album" || album.ArtistID != 10001 {
		t.Error()
	}
}

func TestUpdateScan(t *testing.T) {

	f := []*entities.FileInfo{
		// Existing song with different info
		{
			Album:       "album_2",
			AlbumArtist: "artist_1",
			Path:        "path_2",
			Title:       "Y.M.C.A.",  // Different
			Genre:       "Neurofunk", // Different
		},
		// New song
		{
			Album:       "new_album",
			AlbumArtist: "new_artist",
			Path:        "new_path",
			Title:       "new_title",
		},
	}

	p := provider.NewMockProvider("mock", f)
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{}, m)
	hub := socket.NewMockHub()
	scanner := provider.NewScanner([]provider.Provider{p}, dal, hub)

	// Scan without updating or deleting
	scanner.StartAllScans(true, false)

	// Should change
	s, err := dal.GetSong(2)
	if err != nil || s.Title != "Y.M.C.A." {
		t.Error()
	}

	// New genre should have been added
	genre, err := dal.GetGenre("Neurofunk")
	album, err := dal.GetAlbum(s.AlbumID)
	if err != nil || genre == nil || s.GenreID != genre.ID || album.GenreID != genre.ID {
		t.Error()
	}

	// Song shouldn't change as this song wasn't provided
	s, err = dal.GetSong(1)
	if err != nil {
		t.Error()
	}

	// Song shouldn't change as this song has another provider
	s, err = dal.GetSong(22)
	if err != nil {
		t.Error()
	}

	// New song should have been added
	s, err = dal.GetSong(10001)
	if err != nil || s == nil || s.Title != "new_title" {
		t.Error()
	}

	// New artist should have been added
	artist, err := dal.GetArtist(10001)
	if err != nil || artist == nil || artist.Name != "new_artist" {
		t.Error()
	}

	// New album should have been added
	album, err = dal.GetAlbum(10001)
	if err != nil || album == nil || album.Name != "new_album" || album.ArtistID != 10001 {
		t.Error()
	}

}

func TestDeleteScan(t *testing.T) {

	f := []*entities.FileInfo{
		// Existing song with different info
		{
			Album:       "album_2",
			AlbumArtist: "artist_1",
			Path:        "path_2",
			Title:       "Y.M.C.A.",  // Different
			Genre:       "Neurofunk", // Different
		},
		// New song
		{
			Album:       "new_album",
			AlbumArtist: "artist_1",
			Path:        "new_path",
			Title:       "new_title",
		},
	}

	p := provider.NewMockProvider("mock", f)
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{}, m)
	hub := socket.NewMockHub()
	scanner := provider.NewScanner([]provider.Provider{p}, dal, hub)

	// Scan without updating or deleting
	scanner.StartAllScans(false, true)

	albums := dal.GetAlbums(0, 9999, 0)
	if len(albums) != 2 {
		t.Error()
	}

	artists := dal.Db.GetArtists()
	if len(artists) != 1 {
		t.Error()
	}
}
