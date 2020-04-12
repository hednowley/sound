package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

// TODO: Extract intersection with Artist
type ArtistWithAlbums struct {
	XMLName    xml.Name `xml:"artist" json:"-"`
	ID         uint     `xml:"id,attr" json:"id,string"`
	Name       string   `xml:"name,attr" json:"name"`
	Art        string   `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	AlbumCount uint     `xml:"albumCount,attr" json:"albumCount"`
	Albums     []*Album `xml:"album" json:"album,omitempty"`
	Duration   int      `xml:"duration,attr" json:"duration"`
}

func NewArtistWithAlbums(artist *dao.Artist, albums []dao.Album) *ArtistWithAlbums {

	albumDTOs := make([]*Album, len(albums))
	for index, album := range albums {
		albumDTOs[index] = NewAlbum(&album)
	}

	return &ArtistWithAlbums{
		ID:         artist.ID,
		Name:       artist.Name,
		AlbumCount: artist.AlbumCount,
		Albums:     albumDTOs,
		Art:        artist.Art,
		Duration:   artist.Duration,
	}
}
