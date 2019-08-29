package dao

// Art is a single image of album art.
type Art struct {
	ID   uint   `gorm:"PRIMARY_KEY"`
	Hash string `gorm:"index:art_hash_idx"`
	Path string
}
