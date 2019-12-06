package dal_test

import (
	"sort"
	"testing"

	"github.com/hednowley/sound/dal"

	"github.com/hednowley/sound/dao"
)

func TestPlaylistSorter(t *testing.T) {

	e1 := &dao.PlaylistEntry{
		Index: 4,
	}

	e2 := &dao.PlaylistEntry{
		Index: 0,
	}

	e3 := &dao.PlaylistEntry{
		Index: 10,
	}

	p := dao.Playlist{
		Entries: []*dao.PlaylistEntry{e1, e2, e3},
	}

	sorter := dal.NewPlaylistSorter(&p)
	sort.Sort(sorter)

	if p.Entries[0] != e2 || p.Entries[1] != e1 || p.Entries[2] != e3 {
		t.Error()
	}
}
