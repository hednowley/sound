package dto

import (
	"encoding/xml"
	"path"
	"strings"
	"time"

	"github.com/hednowley/sound/dao"
)

type songBody struct {
	Title       string    `xml:"title,attr" json:"title"`
	AlbumName   string    `xml:"album,attr" json:"album"`
	ArtistName  string    `xml:"artist,attr" json:"artist"`
	Track       int       `xml:"track,attr" json:"track,omitempty"`
	Year        int       `xml:"year,attr,omitempty" json:"year,omitempty"`
	Genre       string    `xml:"genre,attr,omitempty" json:"genre,omitempty"`
	Art         string    `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
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
	ID      uint     `xml:"id,attr" json:"id,string"`
	*songBody
}

func newSongBody(song *dao.Song) *songBody {

	return &songBody{
		AlbumID:    song.AlbumID,
		Title:      song.Title,
		AlbumName:  song.AlbumName,
		ArtistID:   song.AlbumArtistID,
		ArtistName: song.Artist,
		Path:       song.Path,
		Genre:      song.GenreName,
		Year:       song.Year,
		Art:        song.Art,
		Track:      song.Track,
		Disc:       song.Disc,
		Type:       "music",
		IsVideo:    false,
		Created:    *song.Created,
		Extension:  strings.TrimPrefix(path.Ext(song.Path), "."),
		Size:       song.Size,
		Duration:   song.Duration,
		Bitrate:    song.Bitrate,
	}
}

func NewSong(song *dao.Song) *Song {

	return &Song{
		xml.Name{},
		song.ID,
		newSongBody(song),
	}
}
