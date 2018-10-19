package dto

import (
	"encoding/xml"
	"strings"

	"github.com/hednowley/sound/dao"
)

type ArtistCollection struct {
	XMLName         xml.Name       `xml:"artists" json:"-"`
	IgnoredArticles string         `xml:"ignoredArticles,attr" json:"ignoredArticles"`
	Indexes         []*ArtistIndex `xml:"index" json:"index"`
}

type Indexes struct {
	XMLName         xml.Name       `xml:"indexes" json:"-"`
	IgnoredArticles string         `xml:"ignoredArticles,attr" json:"ignoredArticles"`
	Indexes         []*ArtistIndex `xml:"index" json:"index"`
}

type ArtistIndex struct {
	XMLName xml.Name  `xml:"index" json:"-"`
	Name    string    `xml:"name,attr" json:"name"`
	Artists []*Artist `xml:"artist" json:"artist"`
}

func NewArtistIndex(artists []*dao.Artist, name string) *ArtistIndex {

	dtoArtists := make([]*Artist, len(artists))
	for index, a := range artists {
		dtoArtists[index] = NewArtist(a, false)
	}

	return &ArtistIndex{
		Name:    name,
		Artists: dtoArtists,
	}
}

func putArtist(indexes map[string][]*dao.Artist, artist *dao.Artist, letter string) {

	value, exists := indexes[letter]
	if exists {
		value = append(value, artist)
	} else {
		value = []*dao.Artist{artist}
	}

	indexes[letter] = value
}

func makeArtistIndexes(artists []*dao.Artist) []*ArtistIndex {

	indexes := make(map[string][]*dao.Artist)
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for _, a := range artists {

		name := strings.ToUpper(a.Name)
		found := false

		for _, letter := range letters {
			if strings.HasPrefix(name, letter) {
				putArtist(indexes, a, letter)
				found = true
				break
			}
		}

		if !found {
			putArtist(indexes, a, "#")
		}
	}

	o := make([]*ArtistIndex, 27)
	for i, l := range letters {
		o[i] = NewArtistIndex(indexes[l], l)
	}
	o[26] = NewArtistIndex(indexes["#"], "#")

	return o
}

func NewArtistCollection(artists []*dao.Artist) *ArtistCollection {

	indexes := makeArtistIndexes(artists)

	return &ArtistCollection{
		IgnoredArticles: "the",
		Indexes:         indexes,
	}
}

func NewIndexes(artists []*dao.Artist) *Indexes {

	indexes := makeArtistIndexes(artists)

	return &Indexes{
		IgnoredArticles: "the",
		Indexes:         indexes,
	}
}
