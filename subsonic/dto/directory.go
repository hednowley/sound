package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type Directory struct {
	XMLName  xml.Name `xml:"directory" json:"-"`
	ID       string   `xml:"id,attr" json:"id"`
	Name     string   `xml:"name" json:"name"`
	IsDir    bool     `xml:"isDir" json:"isDir"`
	Children []*Song  `xml:"child" json:"child"`
}

func NewAlbumDirectory(album *dao.Album) *Directory {

	songs := make([]*Song, len(album.Songs))
	for index, song := range album.Songs {
		songs[index] = NewSong(song)
	}

	return &Directory{
		ID:       NewAlbumID(album.ID),
		Name:     album.Name,
		IsDir:    true,
		Children: songs,
	}
}

func NewArtistDirectory(artist *dao.Artist) *Directory {

	return &Directory{
		ID:    NewArtistID(artist.ID),
		Name:  artist.Name,
		IsDir: true,
	}
}

func NewSongDirectory(song *dao.Song) *Directory {

	return &Directory{
		ID:    NewSongID(song.ID),
		Name:  song.Title,
		IsDir: false,
	}
}
