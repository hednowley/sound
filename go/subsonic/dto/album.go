package dto

import (
	"encoding/xml"
	"time"

	"github.com/hednowley/sound/dao"
)

type albumBody struct {
	Name      string     `xml:"name,attr" json:"name"`
	Artist    string     `xml:"artist,attr" json:"artist"`
	ArtistID  uint       `xml:"artistId,attr" json:"artistId,string"`
	Art       string     `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	SongCount uint       `xml:"songCount,attr" json:"songCount"`
	Duration  int        `xml:"duration,attr" json:"duration"`
	Created   *time.Time `xml:"created,attr" json:"created"`
	Year      int        `xml:"year,attr,omitempty" json:"year,omitempty"`
	Genre     string     `xml:"genre,attr,omitempty" json:"genre,omitempty"`
}

func newAlbumBody(album *dao.Album) *albumBody {

	return &albumBody{
		Name:      album.Name,
		ArtistID:  album.ArtistID,
		SongCount: album.SongCount,
		Artist:    album.ArtistName,
		Art:       album.Art,
		Created:   album.Created,
		Genre:     album.GenreName,
		Year:      album.Year,
		Duration:  album.Duration,
	}

}

type Album struct {
	XMLName xml.Name `xml:"album" json:"-"`
	ID      uint     `xml:"id,attr" json:"id,string"`
	*albumBody
}

func NewAlbum(album *dao.Album) *Album {
	return &Album{
		XMLName:   xml.Name{},
		ID:        album.ID,
		albumBody: newAlbumBody(album),
	}
}
