package dao

type Genre struct {
	ID     uint `gorm:"PRIMARY_KEY"`
	Name   string
	Songs  []*Song  `gorm:"foreignkey:GenreID"`
	Albums []*Album `gorm:"foreignkey:GenreID"`
}
