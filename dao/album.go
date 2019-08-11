package dao

import (
	"time"
)

type Album struct {
	ID            uint `gorm:"PRIMARY_KEY"`
	Artist        *Artist
	ArtistID      uint    `gorm:"index:albums_name_idx"`
	Name          string  `gorm:"index:albums_name_idx"`
	Songs         []*Song `gorm:"foreignkey:AlbumID"`
	Created       *time.Time
	Art           string
	GenreID       uint
	Genre         *Genre
	Year          int
	Duration      int    // Duration in seconds
	Disambiguator string // Two albums are only considered the same if their Name, Artist and Disambiguator are the same.
	Starred       bool
}
