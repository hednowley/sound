package dto

import (
	"encoding/xml"
	"time"

	"github.com/hednowley/sound/dao"
)

type Album struct {
	XMLName   xml.Name   `xml:"album" json:"-"`
	ID        uint       `xml:"id,attr" json:"id,string"`
	Name      string     `xml:"name,attr" json:"name"`
	Artist    string     `xml:"artist,attr" json:"artist"`
	ArtistID  uint       `xml:"artistId,attr" json:"artistId,string"`
	Art       string     `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
	SongCount int        `xml:"songCount,attr" json:"songCount"`
	Duration  int        `xml:"duration,attr" json:"duration"`
	Created   *time.Time `xml:"created,attr" json:"created"`
	Year      int        `xml:"year,attr,omitempty" json:"year,omitempty"`
	Genre     string     `xml:"genre,attr,omitempty" json:"genre,omitempty"`
	Songs     []*Song    `xml:"song" json:"song,omitempty"`
}

func NewAlbum(album *dao.Album, includeSongs bool) *Album {

	count := len(album.Songs)
	var songs []*Song
	if includeSongs {
		songs = make([]*Song, count)
		for index, song := range album.Songs {
			songs[index] = NewSong(song)
		}
	}

	var artistName string
	if album.Artist != nil {
		artistName = album.Artist.Name
	}

	var genreName string
	if album.Genre != nil {
		genreName = album.Genre.Name
	}

	return &Album{
		ID:        album.ID,
		Name:      album.Name,
		ArtistID:  album.ArtistID,
		SongCount: len(album.Songs),
		Artist:    artistName,
		Songs:     songs,
		Art:       album.Art,
		Created:   album.Created,
		Genre:     genreName,
		Year:      album.Year,
		Duration:  album.Duration,
	}
}
