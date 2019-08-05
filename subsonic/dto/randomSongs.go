package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
)

type RandomSongs struct {
	XMLName xml.Name `xml:"randomSongs" json:"-"`
	Songs   []*Song  `xml:"song" json:"song"`
}

func NewRandomSongs(songs []*dao.Song) *RandomSongs {

	dto := make([]*Song, len(songs))

	for i, s := range songs {
		dto[i] = NewSong(s)
	}
	return &RandomSongs{
		Songs: dto,
	}
}
