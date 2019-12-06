package dao

// Genre is a genre.
type Genre struct {
	ID     uint     `gorm:"PRIMARY_KEY"`
	Name   string   `gorm:"index:genre_name_idx"`
	Songs  []*Song  `gorm:"foreignkey:GenreID"`
	Albums []*Album `gorm:"foreignkey:GenreID"`
}
