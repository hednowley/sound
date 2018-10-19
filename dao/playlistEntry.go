package dao

type PlaylistEntry struct {
	ID         uint `gorm:"PRIMARY_KEY"`
	PlaylistID uint
	SongID     uint
	Song       *Song
	Index      int
}
