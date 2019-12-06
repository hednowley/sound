package dto

import (
	"time"

	"github.com/hednowley/sound/dao"
)

type Album struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	Artist   string         `json:"artist"`
	ArtistID uint           `json:"artistId,string"`
	Art      string         `json:"coverArt,omitempty"`
	Duration int            `json:"duration"`
	Created  *time.Time     `json:"created"`
	Year     int            `json:"year,omitempty"`
	Genre    string         `json:"genre,omitempty"`
	Songs    []*SongSummary `json:"songs"`
}

func NewAlbum(album *dao.Album) *Album {

	songs := make([]*SongSummary, len(album.Songs))
	for index, song := range album.Songs {
		songs[index] = NewSongSummary(song)
	}

	return &Album{
		Name:     album.Name,
		ID:       album.ID,
		ArtistID: album.ArtistID,
		Artist:   album.ArtistName,
		Art:      album.Art,
		Created:  album.Created,
		Genre:    album.GenreName,
		Year:     album.Year,
		Duration: album.Duration,
		Songs:    songs,
	}
}
