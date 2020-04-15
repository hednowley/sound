package dto

import (
	"time"

	"github.com/hednowley/sound/dao"
)

type AlbumSummary struct {
	ID       uint       `json:"id"`
	Name     string     `json:"name"`
	Artist   string     `json:"artist"`
	ArtistID uint       `json:"artistId,string"`
	Art      string     `json:"coverArt,omitempty"`
	Created  *time.Time `json:"created"`
	Year     int        `json:"year,omitempty"`
	Genre    string     `json:"genre,omitempty"`
}

func NewAlbumSummary(album *dao.Album) *AlbumSummary {
	return &AlbumSummary{
		Name:     album.Name,
		ID:       album.ID,
		ArtistID: album.ArtistID,
		Artist:   album.ArtistName,
		Art:      album.GetArt(),
		Created:  album.Created,
		Genre:    album.GetGenre(),
		Year:     album.GetYear(),
	}
}
