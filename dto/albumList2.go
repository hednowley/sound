package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type AlbumList2 struct {
	XMLName xml.Name `xml:"albumList2" json:"-"`
	Albums  []*Album `xml:"album"  json:"album"`
}

func NewAlbumList2(albums []*dao.Album) *AlbumList2 {

	dtoAlbums := make([]*Album, len(albums))
	for index, a := range albums {
		dtoAlbums[index] = NewAlbum(a, false)
	}

	return &AlbumList2{Albums: dtoAlbums}
}
