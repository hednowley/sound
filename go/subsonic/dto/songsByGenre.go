package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type SongsByGenre struct {
	XMLName xml.Name `xml:"songsByGenre" json:"-"`
	Songs   []*Song  `xml:"song" json:"song"`
}

func NewSongsByGenre(songs []dao.Song) *SongsByGenre {

	songDTOs := make([]*Song, len(songs))

	for i, s := range songs {
		songDTOs[i] = NewSong(&s)
	}
	return &SongsByGenre{
		Songs: songDTOs,
	}
}
