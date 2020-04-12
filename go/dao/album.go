package dao

import (
	"time"
)

// Album is an album.
type Album struct {
	ID       uint   `gorm:"PRIMARY_KEY"`
	ArtistID uint   `gorm:"index:albums_artist_id_idx"`
	Name     string `gorm:"index:albums_name_idx"`

	Created       *time.Time
	Art           string
	GenreID       uint
	Year          int
	Duration      int
	Disambiguator string // Two albums are only considered the same if their Name, Artist and Disambiguator are the same.
	Starred       bool

	// Precalculated fields which are stored for performance
	SongCount  uint
	ArtistName string
	GenreName  string
}
