package dao

import "time"

// Song is a song.
type Song struct {
	ID         uint
	Artist     string // Name of the artist of the song. Can differ from the album's artist.
	AlbumID    uint
	Path       string
	Title      string
	Track      int
	Disc       int
	Year       int
	Art        string
	Created    *time.Time
	Size       int64  // File size in bytes
	Bitrate    int    // Bitrate in kb/s
	Duration   int    // Duration in seconds
	Token      string // An ID unique to this song amongst other songs from its provider
	ProviderID string // THe ID of the provider which supplied this song
	Starred    bool

	// Precalculated fields which are stored for performance
	AlbumName     string
	AlbumArtistID uint
	GenreName     string
}
