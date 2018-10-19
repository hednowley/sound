package dao

type Artist struct {
	ID     uint `gorm:"PRIMARY_KEY"`
	Name   string
	Albums []*Album `gorm:"foreignkey:ArtistID"`
	ArtID  uint
	Art    *Art
}
