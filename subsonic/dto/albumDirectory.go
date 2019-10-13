package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type albumDirectoryBody struct {
	*Directory
	*albumBody
}

type AlbumDirectory struct {
	XMLName xml.Name `xml:"directory" json:"-"`
	*albumDirectoryBody
	Children []*SongChildDirectory `xml:"child" json:"child"`
}

type AlbumChildDirectory struct {
	XMLName xml.Name `xml:"child" json:"-"`
	*albumDirectoryBody
}

func newAlbumDirectoryBody(album *dao.Album) *albumDirectoryBody {
	return &albumDirectoryBody{
		Directory: &Directory{
			ID:     NewAlbumID(album.ID),
			IsDir:  true,
			Parent: NewArtistID(album.ArtistID),
		},
		albumBody: newAlbumBody(album, false),
	}
}

func NewAlbumDirectory(album *dao.Album) *AlbumDirectory {
	songs := make([]*SongChildDirectory, len(album.Songs))
	for index, song := range album.Songs {
		songs[index] = NewSongChildDirectory(song)
	}

	return &AlbumDirectory{xml.Name{}, newAlbumDirectoryBody(album), songs}
}

func NewAlbumChildDirectory(album *dao.Album) *AlbumChildDirectory {
	return &AlbumChildDirectory{xml.Name{}, newAlbumDirectoryBody(album)}
}