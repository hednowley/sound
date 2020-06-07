package handler

import (
	"net/url"

	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/dto"
)

func NewGetMusicFoldersHandler(providers []provider.Provider) api.Handler {
	return func(params url.Values, _ *api.HandlerContext) *api.Response {
		return api.NewSuccessfulReponse(dto.NewMusicFolderCollection(providers))
	}
}
