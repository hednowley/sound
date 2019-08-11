package dto

import (
	"errors"
	"strings"

	"github.com/hednowley/sound/util"
)

type DirectoryType int

const (
	SongDirectory   DirectoryType = 0
	AlbumDirectory  DirectoryType = 1
	ArtistDirectory DirectoryType = 2
	songPrefix                    = "song_"
	albumPrefix                   = "album_"
	artistPrefix                  = "artist_"
)

type DirectoryID struct {
	ID   uint
	Type DirectoryType
}

func ParseDirectoryID(s string) (*DirectoryID, error) {
	if strings.HasPrefix(s, songPrefix) {
		id := util.ParseUint(strings.TrimPrefix(s, songPrefix), 0)
		if id == 0 {
			return nil, errors.New("Bad song ID")
		}
		return &DirectoryID{id, SongDirectory}, nil
	}

	if strings.HasPrefix(s, albumPrefix) {
		id := util.ParseUint(strings.TrimPrefix(s, albumPrefix), 0)
		if id == 0 {
			return nil, errors.New("Bad album ID")
		}
		return &DirectoryID{id, AlbumDirectory}, nil
	}

	if strings.HasPrefix(s, artistPrefix) {
		id := util.ParseUint(strings.TrimPrefix(s, artistPrefix), 0)
		if id == 0 {
			return nil, errors.New("Bad artist ID")
		}
		return &DirectoryID{id, ArtistDirectory}, nil
	}

	return nil, errors.New("Unknown directory ID")
}

func NewSongID(id uint) string {
	return songPrefix + util.FormatUint(id)
}

func NewAlbumID(id uint) string {
	return albumPrefix + util.FormatUint(id)
}

func NewArtistID(id uint) string {
	return artistPrefix + util.FormatUint(id)
}
