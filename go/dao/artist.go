package dao

// Artist is an artist.
type Artist struct {
	ID       uint     `gorm:"PRIMARY_KEY"`
	Name     string   `gorm:"index:artists_name_idx"`
	Albums   []*Album `gorm:"foreignkey:ArtistID"`
	Art      string
	Duration int
	Starred  bool

	// Precalculated fields which are stored for performance
	AlbumCount uint
}
