package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type Artist struct {
	XMLName    xml.Name `xml:"artist" json:"-"`
	ID         uint     `xml:"id,attr" json:"id,string"`
	Name       string   `xml:"name,attr" json:"name"`
	Art        string   `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	AlbumCount uint     `xml:"albumCount,attr" json:"albumCount"`
	Albums     []*Album `xml:"album" json:"album,omitempty"`
	Duration   int      `xml:"duration,attr" json:"duration"`
}

func NewArtist(artist *dao.Artist, includeAlbums bool) *Artist {

	var albums []*Album

	if includeAlbums {
		albums = make([]*Album, len(artist.Albums))
		for index, album := range artist.Albums {
			albums[index] = NewAlbum(album, false)
		}
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
