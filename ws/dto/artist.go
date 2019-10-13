package dto

import (
	"github.com/hednowley/sound/dao"
)

type Artist struct {
	ID         uint     `json:"id,string"`
	Name       string   `json:"name"`
	Art        string   `json:"coverArt,omitempty"`
	AlbumCount uint     `json:"albumCount"`
	Albums     []*Album `json:"album,omitempty"`
	Duration   int      `json:"duration"`
}

func NewArtist(artist *dao.Artist) *Artist {

	albums := make([]*Album, len(artist.Albums))
	for index, album := range artist.Albums {
		albums[index] = NewAlbum(album)
	}

	return &Artist{
		ID:         artist.ID,
		Name:       artist.Name,
		AlbumCount: artist.AlbumCount,
		Albums:     albums,
		Art:        artist.Art,
		Duration:   artist.Duration,
	}
}
