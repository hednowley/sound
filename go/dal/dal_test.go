package dal_test

import (
	"testing"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/database"
)

func TestPutPlaylist(t *testing.T) {
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{}, m)

	id, err := dal.PutPlaylist(0, "playlist2", []uint{1, 2, 10}, true)
	if err != nil {
		t.Error(err)
	} else if id != 10001 {
		t.Error("Bad ID returned by PutPlaylist")
	}

	p, err := dal.Db.GetPlaylist(1)
	if err != nil {
		t.Error(err)
	} else if p.EntryCount != 4 {
		t.Error("Wrong entry count")
	}
}

func TestSearchArtist(t *testing.T) {
	m := database.NewMock()

	dal := dal.NewDAL(&config.Config{}, m)

	artists, err := dal.Db.SearchArtists("EethOVE", 10, 0)
	if err != nil {
		t.Error(err)
	} else if len(artists) != 1 {
		t.Error("Could not find beethoven")
	}

	artists, err = dal.Db.SearchArtists("mc hammer", 10, 0)
	if err != nil {
		t.Error(err)
	} else if len(artists) != 0 {
		t.Error("Search returned false positive")
	}

	artists, err = dal.Db.SearchArtists("artist", 2, 0)
	if err != nil {
		t.Error(err)
	} else if len(artists) != 2 {
		t.Errorf("Search is not limiting result count (returned %d)", len(artists))
	}

	artists, err = dal.Db.SearchArtists("artist", 2, 2)
	if err != nil {
		t.Error(err)
	} else if len(artists) != 1 {
		t.Error("Search offset is not working")
	}
}

// func TestPutArt(t *testing.T) {
// 	m := database.NewMock()
// 	dal := dal.NewDAL(&config.Config{}, m)

// 	data1 := &entities.CoverArtData{
// 		Extension: "jpg",
// 		Raw:       []byte("Hello"),
// 	}

// 	a1, err := dal.PutArt(data1)
// 	t.Log(err)
// 	if err != nil || a1 == nil || a1.ID == 0 || a1.Path == "" || a1.Hash == "" {
// 		t.Error("Could not add new art")
// 		return
// 	}

// 	data2 := &entities.CoverArtData{
// 		Extension: "jpg",
// 		Raw:       []byte("Hello"),
// 	}

// 	a2, err := dal.PutArt(data2)
// 	if err != nil || a2 == nil {
// 		t.Error("Could not add new art")
// 		return
// 	}

// 	if a2.ID != a1.ID {
// 		t.Error("Art should be the same")
// 	}

// }
