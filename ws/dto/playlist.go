package dto

import (
	"github.com/hednowley/sound/dao"
)

type Playlist struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"` 
	Songs    []*SongSummary `json:"songs"`
}

func NewPlaylist(playlist *dao.Playlist) *Playlist {

	songs := make([]*SongSummary, len(playlist.Entries))
	for index, entry := range playlist.Entries {
		songs[index] = NewSongSummary(entry.Song)
	}

	return &Playlist{
		Name:     playlist.Name,
		ID:       playlist.ID,
		Songs:    songs,
	}
}
