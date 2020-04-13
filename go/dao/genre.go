package dao

// Genre is a genre.
type Genre struct {
	ID         uint   `gorm:"PRIMARY_KEY"`
	Name       string `gorm:"index:genre_name_idx"`
	SongCount  int
	AlbumCount int
}
