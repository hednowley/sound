package dto

import (
	"encoding/xml"
	"strings"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dao"
)

type ArtistCollection struct {
	XMLName         xml.Name       `xml:"artists" json:"-"`
	IgnoredArticles string         `xml:"ignoredArticles,attr" json:"ignoredArticles"`
	Indexes         []*artistIndex `xml:"index" json:"index"`
}

type indexCollection struct {
	XMLName         xml.Name       `xml:"indexes" json:"-"`
	IgnoredArticles string         `xml:"ignoredArticles,attr" json:"ignoredArticles"`
	Indexes         []*artistIndex `xml:"index" json:"index"`
}

type artistIndex struct {
	XMLName xml.Name  `xml:"index" json:"-"`
	Name    string    `xml:"name,attr" json:"name"`
	Artists []*Artist `xml:"artist" json:"artist"`
}

func newArtistIndex(artists []*dao.Artist, name string) *artistIndex {

	dtoArtists := make([]*Artist, len(artists))
	for index, a := range artists {
		dtoArtists[index] = NewArtist(a, false)
	}

	return &artistIndex{
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

func makeArtistIndexes(artists []*dao.Artist, ignoredArticles []string) ([]*artistIndex, string) {

	indexes := make(map[string][]*dao.Artist)
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for _, a := range artists {

		name := strings.ToUpper(a.Name)
		for _, ia := range ignoredArticles {
			name = strings.TrimPrefix(name, strings.ToUpper(ia)+" ")
		}

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

	o := make([]*artistIndex, 27)
	for i, l := range letters {
		o[i] = newArtistIndex(indexes[l], l)
	}
	o[26] = newArtistIndex(indexes["#"], "#")

	return o, strings.Join(ignoredArticles, " ")
}

func NewArtistCollection(artists []*dao.Artist, conf *config.Config) *ArtistCollection {

	indexes, ignored := makeArtistIndexes(artists, conf.IgnoredArticles)

	return &ArtistCollection{
		IgnoredArticles: ignored,
		Indexes:         indexes,
	}
}

func NewIndexCollection(artists []*dao.Artist, conf *config.Config) *indexCollection {

	indexes, ignored := makeArtistIndexes(artists, conf.IgnoredArticles)

	return &indexCollection{
		IgnoredArticles: ignored,
		Indexes:         indexes,
	}
}
