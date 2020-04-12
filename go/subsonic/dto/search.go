package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type SearchCore struct {
	Artists []*Artist `xml:"artist" json:"artist"`
	Albums  []*Album  `xml:"album" json:"album"`
	Songs   []*Song   `xml:"song" json:"song"`
}

type Search2Response struct {
	XMLName xml.Name `xml:"searchResult2" json:"-"`
	*SearchCore
}

type Search3Response struct {
	XMLName xml.Name `xml:"searchResult3" json:"-"`
	*SearchCore
}

func newSearchResponse(artists []*dao.Artist, albums []*dao.Album, songs []*dao.Song) *SearchCore {

	artistCount := len(artists)
	artistsDto := make([]*Artist, artistCount)

	for i, a := range artists {
		artistsDto[i] = NewArtist(a)
	}

	albumCount := len(albums)
	albumsDto := make([]*Album, albumCount)

	for i, a := range albums {
		albumsDto[i] = NewAlbum(a)
	}

	songCount := len(songs)
	songsDto := make([]*Song, songCount)

	for i, a := range songs {
		songsDto[i] = NewSong(a)
	}

	return &SearchCore{
		artistsDto,
		albumsDto,
		songsDto,
	}
}

func NewSearch2Response(artists []*dao.Artist, albums []*dao.Album, songs []*dao.Song) *Search2Response {
	core := newSearchResponse(artists, albums, songs)
	return &Search2Response{
		xml.Name{},
		core,
	}
}

func NewSearch3Response(artists []*dao.Artist, albums []*dao.Album, songs []*dao.Song) *Search3Response {
	core := newSearchResponse(artists, albums, songs)
	return &Search3Response{
		xml.Name{},
		core,
	}
}
