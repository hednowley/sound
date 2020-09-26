package provider_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/projectpath"
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

	conn, _ := dal.Db.GetConn()
	defer conn.Release()

	// Title shouldn't change
	s, err := dal.Db.GetSong(conn, 2)
	if err != nil {
		t.Error(err)
	} else if s.Title != "title_2" {
		t.Error()
	}

	// Song shouldn't change as this song wasn't provided
	s, err = dal.Db.GetSong(conn, 1)
	if err != nil {
		t.Error(err)
	}

	// Song shouldn't change as this song has another provider
	s, err = dal.Db.GetSong(conn, 22)
	if err != nil {
		t.Error(err)
	}

	// New song should have been added
	s, err = dal.Db.GetSong(conn, 10001)
	if err != nil {
		t.Error(err)
	} else if s == nil || s.Title != "new_title" {
		t.Error()
	}

	// New artist should have been added
	artist, err := dal.Db.GetArtist(conn, 10001)
	if err != nil {
		t.Error(err)
	} else if artist == nil || artist.Name != "new_artist" {
		t.Error()
	}

	// New album should have been added
	album, err := dal.Db.GetAlbum(conn, 10001)
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

	conn, _ := dal.Db.GetConn()
	defer conn.Release()

	// Should change
	s, err := dal.Db.GetSong(conn, 2)
	if err != nil {
		t.Error(err)
	} else if s.Title != "Y.M.C.A." {
		t.Error("Song title not updating")
	}

	// New genre should have been added
	album, err := dal.Db.GetAlbum(conn, s.AlbumID)
	if err != nil {
		t.Error(err)
	} else if !util.ContainsString(album.Genres, "Neurofunk") {
		t.Error("Song genre not propagated to album")
	}

	// Song shouldn't change as this song wasn't provided
	s, err = dal.Db.GetSong(conn, 1)
	if err != nil {
		t.Error(err)
	}

	// Song shouldn't change as this song has another provider
	s, err = dal.Db.GetSong(conn, 22)
	if err != nil {
		t.Error(err)
	}

	// New song should have been added
	s, err = dal.Db.GetSong(conn, 10001)
	if err != nil {
		t.Error(err)
	} else if s == nil {
		t.Error("Song not added")
	} else if s.Title != "new_title" {
		t.Error("Song has wrong title")
	}

	// New artist should have been added
	artist, err := dal.Db.GetArtist(conn, 10001)
	if err != nil {
		t.Error(err)
	} else if artist == nil || artist.Name != "new_artist" {
		t.Error()
	}

	// New album should have been added
	album, err = dal.Db.GetAlbum(conn, 10001)
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

	scanner.StartAllScans(false, true)

	conn, _ := dal.Db.GetConn()
	defer conn.Release()

	albums, err := dal.Db.GetAlbums(conn, dao.AlphabeticalByName, 9999, 0)
	if err != nil {
		t.Error(err)
	} else if len(albums) != 2 {
		t.Errorf("Haven't removed albums (have %v)", len(albums))
	}

	artists, err := dal.Db.GetArtists(conn)
	if err != nil {
		t.Error(err)
	} else if len(artists) != 1 {
		t.Errorf("Haven't removed artists (have %v)", len(artists))
	}

	playlistEntries, err := dal.Db.GetPlaylistSongIds(conn, 1, "tommy")
	if err != nil {
		t.Error(err)
	} else if len(playlistEntries) != 1 {
		t.Errorf("Wrong number of playlist entries (have %v)", len(playlistEntries))
	} else if playlistEntries[0] != 2 {
		t.Error("Have retained wrong playlist entry")
	}
}

func TestBigScan(t *testing.T) {

	tempDir, err := ioutil.TempDir(path.Join(projectpath.Root, "testdata"), "")
	if err != nil {
		t.Error(err)
		return
	}

	defer os.RemoveAll(tempDir)
	defer os.Remove(tempDir)

	art := &entities.CoverArtData{
		Extension: "png",
		Raw:       []byte("Hello"),
	}

	files := []*entities.FileInfo{}
	for i := 0; i < 100; i++ {
		files = append(files, &entities.FileInfo{
			Album:       "big_scan_album",
			Artist:      fmt.Sprintf("big_scan_artist ft. %v", i),
			AlbumArtist: "big_scan_artist",
			Path:        fmt.Sprintf("path_%v", i),
			Title:       fmt.Sprintf("song_%v", i),
			Genre:       fmt.Sprintf("genre_%v", i),
			CoverArt:    art,
		})
	}

	p := provider.NewMockProvider("mock", files)
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{
		ArtPath: tempDir,
	}, m)
	hub := socket.NewMockHub()
	scanner := provider.NewScanner([]provider.Provider{p}, dal, hub)

	// Scan without updating or deleting
	scanner.StartAllScans(false, false)
}
