package dal

import (
	"github.com/hednowley/sound/dao"
)

// PlaylistSorter sorts a playlist's entries by index.
// Should be passed to Go's native sorter.
type PlaylistSorter struct {
	playlist *dao.Playlist
}

func NewPlaylistSorter(playlist *dao.Playlist) *PlaylistSorter {
	return &PlaylistSorter{playlist: playlist}
}

func (s *PlaylistSorter) Len() int {
	return len(s.playlist.Entries)
}

func (s PlaylistSorter) Swap(i, j int) {
	s.playlist.Entries[i], s.playlist.Entries[j] = s.playlist.Entries[j], s.playlist.Entries[i]
}

func (s PlaylistSorter) Less(i, j int) bool {
	return s.playlist.Entries[i].Index < s.playlist.Entries[j].Index
}
