package dao

// Artist is an artist.
type Artist struct {
	ID      uint
	Name    string
	Arts    []string
	Starred bool

	Duration   int
	AlbumCount uint
}

func (a *Artist) GetArt() string {
	if len(a.Arts) > 0 {
		return a.Arts[0]
	}

	return ""
}
