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
	XMLName         xml.Name          `xml:"indexes" json:"-"`
	IgnoredArticles string            `xml:"ignoredArticles,attr" json:"ignoredArticles"`
	Indexes         []*directoryIndex `xml:"index" json:"index"`
}

type artistIndex struct {
	XMLName xml.Name  `xml:"index" json:"-"`
	Name    string    `xml:"name,attr" json:"name"`
	Artists []*Artist `xml:"artist" json:"artist"`
}

type directoryIndex struct {
	XMLName    xml.Name     `xml:"index" json:"-"`
	Name       string       `xml:"name,attr" json:"name"`
	Directorys []*Directory `xml:"artist" json:"artist"`
}

func newArtistIndex(artists []*dao.Artist, letter rune) *artistIndex {

	dtoArtists := make([]*Artist, len(artists))
	for index, a := range artists {
		dtoArtists[index] = NewArtist(a, false)
	}

	return &artistIndex{
		Name:    string(letter),
		Artists: dtoArtists,
	}
}

func newDirectoryIndex(artists []*dao.Artist, letter rune) *directoryIndex {

	dirs := make([]*Directory, len(artists))
	for index, a := range artists {
		dirs[index] = NewArtistDirectory(a)
	}

	return &directoryIndex{
		Name:       string(letter),
		Directorys: dirs,
	}
}

var letters = [...]rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func makeRuneMap(artists []*dao.Artist, ignoredArticles []string) map[rune][]*dao.Artist {

	indexes := map[rune][]*dao.Artist{
		'A': []*dao.Artist{},
		'B': []*dao.Artist{},
		'C': []*dao.Artist{},
		'D': []*dao.Artist{},
		'E': []*dao.Artist{},
		'F': []*dao.Artist{},
		'G': []*dao.Artist{},
		'H': []*dao.Artist{},
		'I': []*dao.Artist{},
		'J': []*dao.Artist{},
		'K': []*dao.Artist{},
		'L': []*dao.Artist{},
		'M': []*dao.Artist{},
		'N': []*dao.Artist{},
		'O': []*dao.Artist{},
		'P': []*dao.Artist{},
		'Q': []*dao.Artist{},
		'R': []*dao.Artist{},
		'S': []*dao.Artist{},
		'T': []*dao.Artist{},
		'U': []*dao.Artist{},
		'V': []*dao.Artist{},
		'W': []*dao.Artist{},
		'X': []*dao.Artist{},
		'Y': []*dao.Artist{},
		'Z': []*dao.Artist{},
		'#': []*dao.Artist{},
	}

	for _, a := range artists {

		name := strings.ToUpper(a.Name)
		for _, ia := range ignoredArticles {
			name = strings.TrimPrefix(name, strings.ToUpper(ia)+" ")
		}

		// The artist with no name
		if len(name) == 0 {
			indexes['#'] = append(indexes['#'], a)
			continue
		}

		found := false
		first := rune(name[0])

		for _, letter := range letters {
			if first == letter {
				indexes[letter] = append(indexes[letter], a)
				found = true
				break
			}
		}

		if !found {
			indexes['#'] = append(indexes['#'], a)
		}
	}

	return indexes
}

func NewArtistCollection(artists []*dao.Artist, conf *config.Config) *ArtistCollection {

	runes := makeRuneMap(artists, conf.IgnoredArticles)

	indexes := make([]*artistIndex, 27)
	for i, l := range letters {
		indexes[i] = newArtistIndex(runes[l], l)
	}
	indexes[26] = newArtistIndex(runes['#'], '#')

	return &ArtistCollection{
		IgnoredArticles: strings.Join(conf.IgnoredArticles, " "),
		Indexes:         indexes,
	}
}

func NewIndexCollection(artists []*dao.Artist, conf *config.Config) *indexCollection {

	runes := makeRuneMap(artists, conf.IgnoredArticles)

	dirs := make([]*directoryIndex, 27)
	for i, l := range letters {
		dirs[i] = newDirectoryIndex(runes[l], l)
	}
	dirs[26] = newDirectoryIndex(runes['#'], '#')

	return &indexCollection{
		IgnoredArticles: strings.Join(conf.IgnoredArticles, " "),
		Indexes:         dirs,
	}
}
