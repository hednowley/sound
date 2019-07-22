package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

type SearchVersion int

const (
	Search2 SearchVersion = 2
	Search3 SearchVersion = 3
)

// NewSearchHandler is a handler for searching for albums, artists and songs.
func NewSearchHandler(dal interfaces.DAL, version SearchVersion) api.Handler {

	return func(params url.Values) *api.Response {

		query := params.Get("query")
		if len(query) == 0 {
			return api.NewErrorReponse(dto.Generic, "No query provided.")
		}

		artistCount := api.ParseUint(params.Get("artistCount"), 20)
		albumCount := api.ParseUint(params.Get("albumCount"), 20)
		songCount := api.ParseUint(params.Get("songCount"), 20)

		artistOffset := api.ParseUint(params.Get("artistOffset"), 0)
		albumOffset := api.ParseUint(params.Get("albumOffset"), 0)
		songOffset := api.ParseUint(params.Get("songOffset"), 0)

		artists := dal.SearchArtists(query, artistCount, artistOffset)
		albums := dal.SearchAlbums(query, albumCount, albumOffset)
		songs := dal.SearchSongs(query, songCount, songOffset)

		var response interface{}
		if version == Search2 {
			response = dto.NewSearch2Response(artists, albums, songs)
		} else {
			response = dto.NewSearch3Response(artists, albums, songs)
		}

		return api.NewSuccessfulReponse(response)
	}
}
