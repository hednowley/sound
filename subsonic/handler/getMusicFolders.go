package handler

import (
	"net/url"

	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetMusicFoldersHandler() api.Handler {
	return func(params url.Values) *api.Response {
		return api.NewSuccessfulReponse(dto.NewMusicFolderCollection(0, "Music"))
	}
}
