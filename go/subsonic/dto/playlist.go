package dto

import (
	"encoding/xml"
	"time"

	"github.com/hednowley/sound/dao"
)

type PlaylistCore struct {
	XMLName   xml.Name  `xml:"playlist" json:"-"`
	ID        uint      `xml:"id,attr" json:"id,string"`
	Name      string    `xml:"name,attr" json:"name"`
	Comment   string    `xml:"comment,attr" json:"comment"`
	Owner     string    `xml:"owner,attr" json:"owner"`
	Public    bool      `xml:"public,attr" json:"public"`
	SongCount int       `xml:"songCount,attr" json:"songCount"`
	Duration  int       `xml:"duration,attr" json:"duration"`
	Created   time.Time `xml:"created,attr" json:"created"`
	Changed   time.Time `xml:"changed,attr" json:"changed"`
	Art       string    `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty"`
}

type Playlist struct {
	XMLName xml.Name `xml:"playlist" json:"-"`
	*PlaylistCore
	Songs []*PlaylistEntry `xml:"entry" json:"entry"`
}

type PlaylistEntry struct {
	XMLName xml.Name `xml:"entry" json:"-"`
	ID      uint     `xml:"id,attr" json:"id,string"`
	*songBody
}

func newPlaylistEntry(song *dao.Song) *PlaylistEntry {
	return &PlaylistEntry{
		xml.Name{},
		song.ID,
		newSongBody(song),
	}
}

func newPlaylistCore(playlist *dao.Playlist) *PlaylistCore {
	return &PlaylistCore{
		ID:        playlist.ID,
		SongCount: playlist.EntryCount,
		Name:      playlist.Name,
		Public:    playlist.Public,
		Changed:   *playlist.Changed,
		Created:   *playlist.Created,
		Comment:   playlist.Comment,
		Duration:  playlist.Duration,
		Owner:     playlist.Owner,
	}
}

func NewPlaylist(playlist *dao.Playlist, playlistSongs []dao.Song) *Playlist {

	count := len(playlistSongs)
	songs := make([]*PlaylistEntry, count)

	for i, s := range playlistSongs {
		songs[i] = newPlaylistEntry(&s)
	}

	return &Playlist{
		xml.Name{},
		newPlaylistCore(playlist),
		songs,
	}
}
