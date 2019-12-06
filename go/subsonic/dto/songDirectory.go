package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type songDirectoryBody struct {
	*Directory
	*songBody
	Name string `xml:"name,attr" json:"name"`
}

type SongDirectory struct {
	XMLName xml.Name `xml:"directory" json:"-"`
	*songDirectoryBody
}

type SongChildDirectory struct {
	XMLName xml.Name `xml:"child" json:"-"`
	*songDirectoryBody
}

func newSongDirectoryBody(song *dao.Song) *songDirectoryBody {
	return &songDirectoryBody{
		Directory: &Directory{
			ID:     NewSongID(song.ID),
			IsDir:  false,
			Parent: NewAlbumID(song.AlbumID),
		},
		songBody: newSongBody(song),
		Name:     song.Title,
	}
}

func NewSongDirectory(song *dao.Song) *SongDirectory {
	return &SongDirectory{xml.Name{}, newSongDirectoryBody(song)}
}

func NewSongChildDirectory(song *dao.Song) *SongChildDirectory {
	return &SongChildDirectory{xml.Name{}, newSongDirectoryBody(song)}
}
