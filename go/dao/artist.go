package dao

// Artist is an artist.
type Artist struct {
	ID      uint   `gorm:"PRIMARY_KEY"`
	Name    string `gorm:"index:artists_name_idx"`
	Art     string
	Starred bool

	// Precalculated fields which are stored for performance
	Duration   int
	AlbumCount uint
}
