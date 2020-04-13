package dao

import (
	"time"
)

// Playlist is a playlist.
type Playlist struct {
	ID      uint `gorm:"PRIMARY_KEY"`
	Name    string
	Comment string
	Public  bool
	Created *time.Time
	Changed *time.Time

	EntryCount int
}
