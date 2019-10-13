package dto

import (
	"github.com/hednowley/sound/dao"
)

type Artist struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Art        string   `json:"coverArt,omitempty"`
	Albums     []*Album `json:"album,omitempty"`
}

func NewArtist(artist *dao.Artist) *Artist {

	albums := make([]*Album, len(artist.Albums))
	for index, album := range artist.Albums {
		albums[index] = NewAlbum(album)
	}

	return &Artist{
		ID:         artist.ID,
		Name:       artist.Name,
		Albums:     albums,
		Art:        artist.Art,
	}
}
