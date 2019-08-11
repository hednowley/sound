package dto

import (
	"errors"
	"strings"

	"github.com/hednowley/sound/util"
)

type DirectoryType int

const (
	SongDirectoryType   DirectoryType = 0
	AlbumDirectoryType  DirectoryType = 1
	ArtistDirectoryType DirectoryType = 2
	albumPrefix                       = "album_"
	artistPrefix                      = "artist_"
)

type DirectoryID struct {
	ID   uint
	Type DirectoryType
}

func ParseDirectoryID(s string) (*DirectoryID, error) {
	if strings.HasPrefix(s, albumPrefix) {
		id := util.ParseUint(strings.TrimPrefix(s, albumPrefix), 0)
		if id == 0 {
			return nil, errors.New("Bad album ID")
		}
		return &DirectoryID{id, AlbumDirectoryType}, nil
	}

	if strings.HasPrefix(s, artistPrefix) {
		id := util.ParseUint(strings.TrimPrefix(s, artistPrefix), 0)
		if id == 0 {
			return nil, errors.New("Bad artist ID")
		}
		return &DirectoryID{id, ArtistDirectoryType}, nil
	}

	id := util.ParseUint(s, 0)
	if id == 0 {
		return nil, errors.New("Unknown directory ID")
	}
	return &DirectoryID{id, SongDirectoryType}, nil

}

func NewSongID(id uint) string {
	return util.FormatUint(id)
}

func NewAlbumID(id uint) string {
	return albumPrefix + util.FormatUint(id)
}

func NewArtistID(id uint) string {
	return artistPrefix + util.FormatUint(id)
}
