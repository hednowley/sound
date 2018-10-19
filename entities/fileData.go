package entities

type FileData struct {
	Path        string
	Artist      string
	Album       string
	AlbumArtist string
	Title       string
	Genre       string
	Track       int
	Disc        int
	Year        int
	CoverArt    *CoverArtData
	Size        int64
	Extension   string
}

type CoverArtData struct {
	Extension string
	Raw       []byte
}
