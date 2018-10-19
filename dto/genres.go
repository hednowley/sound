package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type Genres struct {
	XMLName xml.Name `xml:"genres" json:"-"`
	Genres  []*Genre `xml:"genre" json:"genre"`
}

type Genre struct {
	XMLName    xml.Name `xml:"genre" json:"-"`
	SongCount  int      `xml:"songCount,attr" json:"songCount"`
	AlbumCount int      `xml:"albumCount,attr" json:"albumCount"`
	Name       string   `xml:",chardata" json:"value"`
}

func NewGenre(genre *dao.Genre) *Genre {
	return &Genre{
		SongCount:  len(genre.Songs),
		AlbumCount: len(genre.Albums),
		Name:       genre.Name,
	}
}

func NewGenres(genres []*dao.Genre) *Genres {

	count := len(genres)
	dtoGenres := make([]*Genre, count)

	for i, g := range genres {
		dtoGenres[i] = NewGenre(g)
	}

	return &Genres{
		Genres: dtoGenres,
	}
}
