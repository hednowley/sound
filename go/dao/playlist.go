package dao

import (
	"time"
)

// Playlist is a playlist.
type Playlist struct {
	ID       uint
	Name     string
	Comment  string
	Public   bool
	Created  *time.Time
	Changed  *time.Time
	Duration int
	Owner    string

	EntryCount int
}
