package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type ArtistDirectory struct {
	*ArtistDirectorySummary
	Children []*AlbumChildDirectory `xml:"child" json:"child"`
}

type ArtistDirectorySummary struct {
	XMLName xml.Name `xml:"directory" json:"-"`
	*Directory
	Name string `xml:"name,attr" json:"name"`
}

func NewArtistDirectorySummary(artist *dao.Artist) *ArtistDirectorySummary {
	return &ArtistDirectorySummary{
		XMLName: xml.Name{},
		Directory: &Directory{
			ID:    NewArtistID(artist.ID),
			IsDir: true,
		},
		Name: artist.Name,
	}
}

func NewArtistDirectory(artist *dao.Artist) *ArtistDirectory {

	albums := make([]*AlbumChildDirectory, len(artist.Albums))
	for i, a := range artist.Albums {
		albums[i] = NewAlbumChildDirectory(a)
	}

	return &ArtistDirectory{
		NewArtistDirectorySummary(artist),
		albums,
	}
}
