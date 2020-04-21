package provider_test

import (
	"testing"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/util"

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
	s, err := dal.Db.GetSong(2)
	if err != nil {
		t.Error(err)
	} else if s.Title != "title_2" {
		t.Error()
	}

	// Song shouldn't change as this song wasn't provided
	s, err = dal.Db.GetSong(1)
	if err != nil {
		t.Error(err)
	}

	// Song shouldn't change as this song has another provider
	s, err = dal.Db.GetSong(22)
	if err != nil {
		t.Error(err)
	}

	// New song should have been added
	s, err = dal.Db.GetSong(10001)
	if err != nil {
		t.Error(err)
	} else if s == nil || s.Title != "new_title" {
		t.Error()
	}

	// New artist should have been added
	artist, err := dal.Db.GetArtist(10001)
	if err != nil {
		t.Error(err)
	} else if artist == nil || artist.Name != "new_artist" {
		t.Error()
	}

	// New album should have been added
	album, err := dal.Db.GetAlbum(10001)
	if err != nil {
		t.Error(err)
	} else if album == nil || album.Name != "new_album" || album.ArtistID != 10001 {
		t.Error()
	}
}

func TestUpdateScan(t *testing.T) {

	f := []*entities.FileInfo{
		// Existing song with different info
		{
			Album:       "album_1",
			Artist:      "artist_1 ft. Pitbull",
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
	s, err := dal.Db.GetSong(2)
	if err != nil {
		t.Error(err)
	} else if s.Title != "Y.M.C.A." {
		t.Error("Song title not updating")
	}

	// New genre should have been added
	album, err := dal.Db.GetAlbum(s.AlbumID)
	if err != nil {
		t.Error(err)
	} else if !util.ContainsString(album.Genres, "Neurofunk") {
		t.Error("Song genre not propagated to album")
	}

	// Song shouldn't change as this song wasn't provided
	s, err = dal.Db.GetSong(1)
	if err != nil {
		t.Error(err)
	}

	// Song shouldn't change as this song has another provider
	s, err = dal.Db.GetSong(22)
	if err != nil {
		t.Error(err)
	}

	// New song should have been added
	s, err = dal.Db.GetSong(10001)
	if err != nil {
		t.Error(err)
	} else if s == nil {
		t.Error("Song not added")
	} else if s.Title != "new_title" {
		t.Error("Song has wrong title")
	}

	// New artist should have been added
	artist, err := dal.Db.GetArtist(10001)
	if err != nil {
		t.Error(err)
	} else if artist == nil || artist.Name != "new_artist" {
		t.Error()
	}

	// New album should have been added
	album, err = dal.Db.GetAlbum(10001)
	if err != nil {
		t.Error(err)
	} else if album == nil || album.Name != "new_album" || album.ArtistID != 10001 {
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

	albums, err := dal.Db.GetAlbums(dao.AlphabeticalByName, 9999, 0)
	if err != nil {
		t.Error(err)
	} else if len(albums) != 2 {
		t.Errorf("Haven't removed albums (have %v)", len(albums))
	}

	artists, err := dal.Db.GetArtists()
	if err != nil {
		t.Error(err)
	} else if len(artists) != 1 {
		t.Errorf("Haven't removed artists (have %v)", len(artists))
	}
}
