package dto

import "github.com/hednowley/sound/dao"

type AlbumCollection struct {
	Albums []*AlbumSummary `json:"albums"`
}

func NewAlbumCollection(albums []*dao.Album) *AlbumCollection {

	dtoAlbums := make([]*AlbumSummary, len(albums))
	for index, a := range albums {
		dtoAlbums[index] = NewAlbumSummary(a)
	}

	return &AlbumCollection{
		Albums: dtoAlbums,
	}
}
