package entities

type FileInfo struct {
	Path          string
	Artist        string
	Album         string
	AlbumArtist   string
	Title         string
	Genre         string
	Track         int
	Disc          int
	Year          int
	CoverArt      *CoverArtData
	Size          int64
	Extension     string
	Bitrate       int // Bitrate in kb/s
	Duration      int // Duration in seconds
	Disambiguator string
}

type CoverArtData struct {
	Extension string
	Raw       []byte
}
