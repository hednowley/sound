package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type PlaylistCollection struct {
	XMLName   xml.Name                  `xml:"playlists" json:"-"`
	Playlists []*PlaylistCollectionItem `xml:"playlist" json:"playlist"`
}

type PlaylistCollectionItem struct {
	XMLName xml.Name `xml:"playlist" json:"-"`
	*PlaylistCore
}

func NewPlaylistCollectionItem(playlist *dao.Playlist) *PlaylistCollectionItem {
	return &PlaylistCollectionItem{
		xml.Name{},
		newPlaylistCore(playlist),
	}
}

func NewPlaylistCollection(playlists []*dao.Playlist) *PlaylistCollection {

	count := len(playlists)
	dtoCollection := make([]*PlaylistCollectionItem, count)

	for i, p := range playlists {
		dtoCollection[i] = NewPlaylistCollectionItem(p)
	}

	return &PlaylistCollection{
		Playlists: dtoCollection,
	}
}
