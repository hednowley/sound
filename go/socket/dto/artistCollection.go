package dto

import "github.com/hednowley/sound/dao"

type ArtistCollection struct {
	Artists []*ArtistSummary `json:"artists"`
}

type ArtistSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewArtistSummary(artist *dao.Artist) *ArtistSummary {
	return &ArtistSummary{
		ID:   artist.ID,
		Name: artist.Name,
	}
}

func NewArtistCollection(artists []dao.Artist) *ArtistCollection {

	dtoArtists := make([]*ArtistSummary, len(artists))
	for index, a := range artists {
		dtoArtists[index] = NewArtistSummary(&a)
	}

	return &ArtistCollection{
		Artists: dtoArtists,
	}
}
