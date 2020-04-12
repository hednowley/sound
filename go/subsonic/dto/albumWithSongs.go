package dto

import (
	"github.com/hednowley/sound/dao"
)

type albumWithSongsBody struct {
	Songs []*Song `xml:"song" json:"song,omitempty"`
}

func newAlbumWithSongsBody(album *dao.Album, songs []dao.Song) *albumWithSongsBody {

	songsDto := make([]*Song, len(songs))
	for index, song := range songs {
		songsDto[index] = NewSong(&song)

	}

	return &albumWithSongsBody{

		Songs: songsDto,
	}

}

type AlbumWithSongs struct {
	*albumBody
	*albumWithSongsBody
}

func NewAlbumWithSongs(album *dao.Album, songs []dao.Song) *AlbumWithSongs {
	return &AlbumWithSongs{
		albumBody:          newAlbumBody(album),
		albumWithSongsBody: newAlbumWithSongsBody(album, songs),
	}
}
