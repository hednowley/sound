package dal_test

import (
	"testing"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/database"
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
