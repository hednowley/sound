package dto

import (
	"encoding/xml"

	"github.com/hednowley/sound/dao"
	"github.com/hednowley/sound/util"
)

type SongsByGenre struct {
	XMLName xml.Name `xml:"songsByGenre" json:"-"`
	Songs   []*Song  `xml:"song" json:"song"`
}

func NewSongsByGenre(genre *dao.Genre, count uint, offset uint) *SongsByGenre {

	l := len(genre.Songs)
	o := int(offset)
	c := int(count)

	o = util.Max(util.Min(l, o), 0)
	c = util.Min(c, l-o)

	songs := make([]*Song, c)

	for i, s := range genre.Songs[o : o+c] {
		songs[i] = NewSong(s)
	}
	return &SongsByGenre{
		Songs: songs,
	}
}
