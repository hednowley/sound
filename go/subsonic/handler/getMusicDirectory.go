package handler

import (
	"net/url"

	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

// NewGetMusicDirectoryHandler does http://www.subsonic.org/pages/api.jsp#getMusicDirectory
func NewGetMusicDirectoryHandler(dal *dal.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		id, err := dto.ParseDirectoryID(params.Get("id"))
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		switch id.Type {
		case dto.ArtistDirectoryType:
			artist, err := dal.GetArtist(id.ID)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}
			return api.NewSuccessfulReponse(dto.NewArtistDirectory(artist))

		case dto.AlbumDirectoryType:
			album, err := dal.GetAlbum(id.ID, false, false, true)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}
			return api.NewSuccessfulReponse(dto.NewAlbumDirectory(album))

		case dto.SongDirectoryType:
			song, err := dal.GetSong(id.ID, false, false, false)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}
			return api.NewSuccessfulReponse(dto.NewSongDirectory(song))
		}

		return api.NewErrorReponse(dto.Generic, "Unknown ID")
	}
}
