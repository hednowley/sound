package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/provider"
)

type MusicFolderCollection struct {
	XMLName xml.Name       `xml:"musicFolders" json:"-"`
	Folders []*musicFolder `xml:"musicFolder,attr" json:"musicFolder"`
}

type musicFolder struct {
	XMLName xml.Name `xml:"musicFolder" json:"-"`
	ID      int      `xml:"id,attr" json:"id"`
	Name    string   `xml:"name" json:"name"`
}

func NewMusicFolderCollection(providers []provider.Provider) *MusicFolderCollection {

	folders := make([]*musicFolder, len(providers))
	for i, p := range providers {
		folders[i] = &musicFolder{ID: i, Name: p.ID()}
	}

	return &MusicFolderCollection{Folders: folders}
}
