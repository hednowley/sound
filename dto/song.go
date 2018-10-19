package dto

import (
	"encoding/xml"
	"time"

	"github.com/hednowley/sound/dao"
)

type songBody struct {
	ID          uint      `xml:"id,attr" json:"id,string"`
	Parent      uint      `xml:"parent,attr" json:"parent,string"`
	IsDir       bool      `xml:"isDir,attr" json:"isDir"`
	Title       string    `xml:"title,attr" json:"title"`
	AlbumName   string    `xml:"album,attr" json:"album"`
	ArtistName  string    `xml:"artist,attr" json:"artist"`
	Track       int       `xml:"track,attr" json:"track,omitempty"`
	Year        int       `xml:"year,attr,omitempty" json:"year,omitempty"`
	Genre       string    `xml:"genre,attr,omitempty" json:"genre,omitempty"`
	ArtID       uint      `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty,string"`
	Size        int64     `xml:"size,attr" json:"size"`
	ContentType string    `xml:"contentType,attr" json:"contentType"`
	Extension   string    `xml:"suffix,attr" json:"suffix"`
	Duration    int       `xml:"duration,attr" json:"duration"`
	Bitrate     int       `xml:"bitRate,attr" json:"bitRate"`
	Path        string    `xml:"path,attr" json:"path"`
	IsVideo     bool      `xml:"isVideo,attr" json:"isVideo"`
	PlayCount   int       `xml:"playCount,attr" json:"playCount"`
	Disc        int       `xml:"discNumber,attr,omitempty" json:"discNumber,omitempty"`
	Created     time.Time `xml:"created,attr" json:"created"`
	AlbumID     uint      `xml:"albumId,attr" json:"albumId,string"`
	ArtistID    uint      `xml:"artistId,attr" json:"artistId,string"`
	Type        string    `xml:"type,attr" json:"type"`
}

type Song struct {
	XMLName xml.Name `xml:"song" json:"-"`
	*songBody
}

func newSongBody(song *dao.Song) *songBody {

	// Fix for playlists with nonsense songs
	if song == nil {
		song = &dao.Song{}
	}

	var genreName string
	if song.Genre != nil {
		genreName = song.Genre.Name
	}

	return &songBody{
		ID:         song.ID,
		Parent:     song.AlbumID,
		AlbumID:    song.AlbumID,
		Title:      song.Title,
		AlbumName:  song.Album.Name,
		ArtistID:   song.Album.ArtistID,
		ArtistName: song.Album.Artist.Name,
		Path:       song.Path,
		Genre:      genreName,
		Year:       song.Year,
		ArtID:      song.ArtID,
		Track:      song.Track,
		Disc:       song.Disc,
		IsDir:      false,
		Type:       "music",
		IsVideo:    false,
		Created:    *song.Created,
		Extension:  song.Extension,
		Size:       song.Size,
	}
}

func NewSong(song *dao.Song) *Song {

	return &Song{
		xml.Name{},
		newSongBody(song),
	}
}
