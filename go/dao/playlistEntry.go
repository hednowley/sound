package dao

// PlaylistEntry is a single instance of a song inside a playlist.
type PlaylistEntry struct {
	ID         uint
	PlaylistID uint
	SongID     uint
	Index      int
}
