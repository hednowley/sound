package dao

import (
	"time"
)

// Album is an album.
type Album struct {
	ID       uint
	ArtistID uint
	Name     string

	Created       *time.Time
	Arts          []string
	Genres        []string
	Years         []int
	Duration      int
	Disambiguator string // Two albums are only considered the same if their Name, Artist and Disambiguator are the same.
	Starred       bool

	SongCount  uint
	ArtistName string
}

func (a *Album) GetArt() string {
	if len(a.Arts) > 0 {
		return a.Arts[0]
	}

	return ""
}

func (a *Album) GetGenre() string {
	if len(a.Genres) > 0 {
		return a.Genres[0]
	}

	return ""
}

func (a *Album) GetYear() int {
	if len(a.Years) > 0 {
		return a.Years[0]
	}

	return 0
}
