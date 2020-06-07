package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
	"github.com/hednowley/sound/util"
)

type SearchVersion int

const (
	Search2 SearchVersion = 2
	Search3 SearchVersion = 3
)

// NewSearchHandler is a handler for searching for albums, artists and songs.
func NewSearchHandler(dal *dal.DAL, version SearchVersion) api.Handler {

	return func(params url.Values, _ *api.HandlerContext) *api.Response {

		query := params.Get("query")
		if len(query) == 0 {
			return api.NewErrorReponse(dto.Generic, "No query provided.")
		}

		artistCount := util.ParseUint(params.Get("artistCount"), 20)
		albumCount := util.ParseUint(params.Get("albumCount"), 20)
		songCount := util.ParseUint(params.Get("songCount"), 20)

		artistOffset := util.ParseUint(params.Get("artistOffset"), 0)
		albumOffset := util.ParseUint(params.Get("albumOffset"), 0)
		songOffset := util.ParseUint(params.Get("songOffset"), 0)

		conn, err := dal.Db.GetConn()
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}
		defer conn.Release()

		artists, err := dal.Db.SearchArtists(conn, query, artistCount, artistOffset)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		albums, err := dal.Db.SearchAlbums(conn, query, albumCount, albumOffset)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		songs, err := dal.Db.SearchSongs(conn, query, songCount, songOffset)
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		var response interface{}
		if version == Search2 {
			response = dto.NewSearch2Response(artists, albums, songs)
		} else {
			response = dto.NewSearch3Response(artists, albums, songs)
		}

		return api.NewSuccessfulReponse(response)
	}
}
