package handler

import (
	"net/url"

	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetMusicDirectoryHandler(database interfaces.DAL) api.Handler {

	return func(params url.Values) *api.Response {

		id, err := dto.ParseDirectoryID(params.Get("id"))
		if err != nil {
			return api.NewErrorReponse(dto.Generic, err.Error())
		}

		switch id.Type {
		case dto.ArtistDirectory:
			artist, err := database.GetArtist(id.ID)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}
			return api.NewSuccessfulReponse(dto.NewArtistDirectory(artist))

		case dto.AlbumDirectory:
			album, err := database.GetAlbum(id.ID, true, true, true)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}
			return api.NewSuccessfulReponse(dto.NewAlbumDirectory(album))

		case dto.SongDirectory:
			song, err := database.GetSong(id.ID, true, true, true)
			if err != nil {
				return api.NewErrorReponse(dto.Generic, err.Error())
			}
			return api.NewSuccessfulReponse(dto.NewSongDirectory(song))
		}

		return api.NewErrorReponse(dto.Generic, "Unknown ID")
	}
}
