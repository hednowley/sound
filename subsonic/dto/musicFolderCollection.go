package dto

import (
	"encoding/xml"
)

type MusicFolderCollection struct {
	XMLName xml.Name `xml:"directory" json:"-"`
	ID      uint     `xml:"id,attr" json:"id"`
	Name    string   `xml:"name" json:"name"`
}

func NewMusicFolderCollection(id uint, name string) *MusicFolderCollection {
	return &MusicFolderCollection{
		ID:   id,
		Name: name,
	}
}
