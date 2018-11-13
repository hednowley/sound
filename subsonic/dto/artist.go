package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type Artist struct {
	XMLName    xml.Name `xml:"artist" json:"-"`
	ID         uint     `xml:"id,attr" json:"id,string"`
	Name       string   `xml:"name,attr" json:"name"`
	ArtID      uint     `xml:"coverArt,attr,omitempty" json:"coverArt,string,omitempty"`
	AlbumCount int      `xml:"albumCount,attr" json:"albumCount"`
	Albums     []*Album `xml:"album" json:"album,omitempty"`
	Duration   int      `xml:"duration,attr" json:"duration"`
}

func NewArtist(artist *dao.Artist, includeAlbums bool) *Artist {

	albumCount := len(artist.Albums)
	var albums []*Album

	if includeAlbums {
		albums = make([]*Album, albumCount)
		for index, album := range artist.Albums {
			albums[index] = NewAlbum(album, false)
		}
	}

	return &Artist{
		ID:         artist.ID,
		Name:       artist.Name,
		AlbumCount: albumCount,
		Albums:     albums,
		ArtID:      artist.ArtID,
		Duration:   artist.Duration,
	}
}
