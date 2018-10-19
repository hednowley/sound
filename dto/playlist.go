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
	ArtID     uint      `xml:"coverArt,attr,omitempty" json:"coverArt,omitempty,string"`
}

type Playlist struct {
	XMLName xml.Name `xml:"playlist" json:"-"`
	*PlaylistCore
	Songs []*PlaylistEntry `xml:"entry" json:"entry"`
}

type PlaylistEntry struct {
	XMLName xml.Name `xml:"entry" json:"-"`
	*songBody
}

func newPlaylistEntry(song *dao.Song) *PlaylistEntry {
	return &PlaylistEntry{
		xml.Name{},
		newSongBody(song),
	}
}

func newPlaylistCore(playlist *dao.Playlist) *PlaylistCore {
	return &PlaylistCore{
		ID:        playlist.ID,
		SongCount: len(playlist.Entries),
		Name:      playlist.Name,
		Public:    playlist.Public,
		Changed:   *playlist.Changed,
		Created:   *playlist.Created,
		Comment:   playlist.Comment,
		Owner:     "ned",
	}
}

func NewPlaylist(playlist *dao.Playlist) *Playlist {

	count := len(playlist.Entries)
	songs := make([]*PlaylistEntry, count)

	for i, e := range playlist.Entries {
		songs[i] = newPlaylistEntry(e.Song)
	}

	return &Playlist{
		xml.Name{},
		newPlaylistCore(playlist),
		songs,
	}
}
