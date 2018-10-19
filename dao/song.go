package dao

import "time"

type Song struct {
	ID        uint `gorm:"PRIMARY_KEY"`
	Album     *Album
	AlbumID   uint
	Path      string
	Title     string
	Track     int
	Disc      int
	GenreID   uint
	Genre     *Genre
	Year      int
	ArtID     uint
	Art       *Art
	Created   *time.Time
	Extension string
	Size      int64
}
