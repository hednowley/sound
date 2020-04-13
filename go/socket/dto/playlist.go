package dto

import (
	"github.com/hednowley/sound/dao"
)

type Playlist struct {
	ID    uint           `json:"id"`
	Name  string         `json:"name"`
	Songs []*SongSummary `json:"songs"`
}

func NewPlaylist(playlist *dao.Playlist, playlistSongs []dao.Song) *Playlist {

	songs := make([]*SongSummary, len(playlistSongs))
	for index, song := range playlistSongs {
		songs[index] = NewSongSummary(&song)
	}

	return &Playlist{
		Name:  playlist.Name,
		ID:    playlist.ID,
		Songs: songs,
	}
}
