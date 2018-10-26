package dao

import "time"

type Song struct {
	ID         uint `gorm:"PRIMARY_KEY"`
	Album      *Album
	AlbumID    uint
	Path       string `gorm:"index:songs_path_idx"`
	Title      string
	Track      int
	Disc       int
	GenreID    uint
	Genre      *Genre
	Year       int
	ArtID      uint
	Art        *Art
	Created    *time.Time
	Extension  string
	Size       int64
	Bitrate    int // Bitrate in kb/s
	Duration   int // Duration in seconds
	ScanID     string
	ProviderID string `gorm:"index:songs_path_idx"`
}
