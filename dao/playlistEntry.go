package dao

// PlaylistEntry is a single instance of a song inside a playlist.
type PlaylistEntry struct {
	ID         uint `gorm:"PRIMARY_KEY"`
	PlaylistID uint
	SongID     uint
	Song       *Song
	Index      int 
}
