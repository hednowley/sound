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
	Duration   int      `xml:"duration,attr" json:"duration"`
}

func NewArtist(artist *dao.Artist) *Artist {

	return &Artist{
		ID:         artist.ID,
		Name:       artist.Name,
		AlbumCount: artist.AlbumCount,
		Art:        artist.GetArt(),
		Duration:   artist.Duration,
	}
}
