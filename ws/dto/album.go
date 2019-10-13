package dto

import (
	"time"

	"github.com/hednowley/sound/dao"
)

type Album struct {
	ID        uint       `json:"id,string"`
	Name      string     `json:"name"`
	Artist    string     `json:"artist"`
	ArtistID  uint       `json:"artistId,string"`
	Art       string     `json:"coverArt,omitempty"`
	SongCount uint       `json:"songCount"`
	Duration  int        `json:"duration"`
	Created   *time.Time `json:"created"`
	Year      int        `json:"year,omitempty"`
	Genre     string     `json:"genre,omitempty"`
}


func NewAlbum(album *dao.Album) *Album {
	return &Album{
		Name:      album.Name,
		ID:        album.ID,
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