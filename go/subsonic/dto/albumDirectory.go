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
		albumBody: newAlbumBody(album),
	}
}

func NewAlbumDirectory(album *dao.Album, songs []dao.Song) *AlbumDirectory {
	children := make([]*SongChildDirectory, len(songs))
	for index, song := range songs {
		children[index] = NewSongChildDirectory(&song)
	}

	return &AlbumDirectory{xml.Name{}, newAlbumDirectoryBody(album), children}
}

func NewAlbumChildDirectory(album *dao.Album) *AlbumChildDirectory {
	return &AlbumChildDirectory{xml.Name{}, newAlbumDirectoryBody(album)}
}
