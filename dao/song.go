package dao

import "time"

type Song struct {
	ID         uint `gorm:"PRIMARY_KEY"`
	Artist     string
	Album      *Album
	AlbumID    uint
	Path       string
	Title      string
	Track      int
	Disc       int
	GenreID    uint
	Genre      *Genre
	Year       int
	Art        string
	Created    *time.Time
	Extension  string // File extension (without leading full stop)
	Size       int64  // File size in bytes
	Bitrate    int    // Bitrate in kb/s
	Duration   int    // Duration in seconds
	Token      string `gorm:"index:songs_token_idx"` // An ID unique to this song amongst other songs from its provider
	ProviderID string `gorm:"index:songs_token_idx"` // THe ID of the provider which supplied this song
}
