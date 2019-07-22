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
