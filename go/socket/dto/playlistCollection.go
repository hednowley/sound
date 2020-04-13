package dto

import (
	"github.com/hednowley/sound/dao"
)

type playlistSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func newPlaylistSummary(playlist *dao.Playlist) *playlistSummary {
	return &playlistSummary{
		Name: playlist.Name,
		ID:   playlist.ID,
	}
}

type PlaylistCollection struct {
	Playlists []*playlistSummary `json:"playlists"`
}

func NewPlaylistCollection(playlists []dao.Playlist) *PlaylistCollection {
	list := make([]*playlistSummary, len(playlists))
	for index, p := range playlists {
		list[index] = newPlaylistSummary(&p)
	}

	return &PlaylistCollection{
		Playlists: list,
	}
}
