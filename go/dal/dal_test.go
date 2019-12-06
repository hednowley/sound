package dal_test

import (
	"testing"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/entities"
)

func TestPutPlaylist(t *testing.T) {
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{}, m)

	id, err := dal.PutPlaylist(0, "playlist2", []uint{1, 2, 10})
	if err != nil || id != 10001 {
		t.Error()
	}

	p, err := dal.GetPlaylist(1)
	if err != nil {
		t.Error()
	}

	if len(p.Entries) != 4 {
		t.Error()
	}
}

func TestSearchArtist(t *testing.T) {
	m := database.NewMock()

	m.PutArtist(&dao.Artist{
		ID:   100,
		Name: "beethoven",
	})

	dal := dal.NewDAL(&config.Config{}, m)

	artists := dal.SearchArtists("EethOVE", 10, 0)
	if len(artists) != 1 {
		t.Error("Could not find beethoven")
	}

	artists = dal.SearchArtists("mc hammer", 10, 0)
	if len(artists) != 0 {
		t.Error("Search returned false positive")
	}

	artists = dal.SearchArtists("artist", 2, 0)
	if len(artists) != 2 {
		t.Error("Search is not limiting result count")
	}

	artists = dal.SearchArtists("artist", 2, 2)
	if len(artists) != 1 {
		t.Error("Search offset is not working")
	}
}

func TestPutArt(t *testing.T) {
	m := database.NewMock()
	dal := dal.NewDAL(&config.Config{}, m)

	data1 := &entities.CoverArtData{
		Extension: "jpg",
		Raw:       []byte("Hello"),
	}

	a1 := dal.PutArt(data1)
	if a1 == nil || a1.ID == 0 || a1.Path == "" || a1.Hash == "" {
		t.Error("Could not add new art")
	}

	data2 := &entities.CoverArtData{
		Extension: "jpg",
		Raw:       []byte("Hello"),
	}

	a2 := dal.PutArt(data2)
	if a2.ID != a1.ID {
		t.Error("Art should be the same")
	}
}
