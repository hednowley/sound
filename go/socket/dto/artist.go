package dto

import (
	"github.com/hednowley/sound/dao"
)

type Artist struct {
	ID     uint            `json:"id"`
	Name   string          `json:"name"`
	Art    string          `json:"coverArt,omitempty"`
	Albums []*AlbumSummary `json:"albums"`
}

func NewArtist(artist *dao.Artist, albums []dao.Album) *Artist {

	albumsDto := make([]*AlbumSummary, len(albums))
	for index, album := range albums {
		albumsDto[index] = NewAlbumSummary(&album)
	}

	return &Artist{
		ID:     artist.ID,
		Name:   artist.Name,
		Albums: albumsDto,
		Art:    artist.Art,
	}
}
