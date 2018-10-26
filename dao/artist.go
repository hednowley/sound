package dao

type Artist struct {
	ID       uint     `gorm:"PRIMARY_KEY"`
	Name     string   `gorm:"index:artists_name_idx"`
	Albums   []*Album `gorm:"foreignkey:ArtistID"`
	ArtID    uint
	Art      *Art
	Duration int
}
