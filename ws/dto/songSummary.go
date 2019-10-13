package dto

import (
	"github.com/hednowley/sound/dao"
)

type SongSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewSongSummary(song *dao.Song) *SongSummary {
	return &SongSummary{
		Name: song.Title,
		ID:   song.ID,
	}
}
