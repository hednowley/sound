package dao

import (
	"time"
)

type Album struct {
	ID       uint `gorm:"PRIMARY_KEY"`
	Artist   *Artist
	ArtistID uint
	Name     string
	Songs    []*Song `gorm:"foreignkey:AlbumID"`
	Created  *time.Time
	ArtID    uint
	Art      *Art
	GenreID  uint
	Genre    *Genre
	Year     int
}
