package dao

type Art struct {
	ID   uint `gorm:"PRIMARY_KEY"`
	Hash string
	Path string
}
