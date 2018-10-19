package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type Directory struct {
	XMLName  xml.Name `xml:"directory" json:"-"`
	ID       uint     `xml:"id,attr" json:"id"`
	Name     string   `xml:"name" json:"name"`
	IsDir    bool     `xml:"isDir" json:"isDir"`
	Children []*Song  `xml:"child" json:"child"`
}

func NewDirectory(album *dao.Album) Directory {

	songs := make([]*Song, len(album.Songs))
	for index, song := range album.Songs {
		songs[index] = NewSong(song)
	}

	return Directory{
		ID:       album.ID,
		Name:     album.Name,
		IsDir:    true,
		Children: songs,
	}
}
