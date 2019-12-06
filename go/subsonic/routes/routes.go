package routes

import (
	"net/http"

	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/handler"
)

// Routes describes which routes correspond to which HTTP handlers.
type Routes map[string]http.HandlerFunc

// NewRoutes constructs a new description of the routes for the Subsonic API.
func NewRoutes(factory *api.HandlerFactory, config *config.Config, dal interfaces.DAL, scanner *provider.Scanner, providers []provider.Provider) Routes {
	routes := make(map[string]http.HandlerFunc)

	routes["/subsonic/rest/ping"] = factory.PublishHandler(handler.NewPingHandler())
	routes["/subsonic/rest/getlicense"] = factory.PublishHandler(handler.NewGetLicenseHandler())

	// Scanning
	routes["/subsonic/rest/getscanstatus"] = factory.PublishHandler(handler.NewGetScanStatusHandler(scanner))
	routes["/subsonic/rest/startscan"] = factory.PublishHandler(handler.NewStartScanHandler(scanner))

	// Querying
	routes["/subsonic/rest/getalbumlist2"] = factory.PublishHandler(handler.NewGetAlbumList2Handler(dal))
	routes["/subsonic/rest/getartists"] = factory.PublishHandler(handler.NewGetArtistsHandler(dal, config))
	routes["/subsonic/rest/getindexes"] = factory.PublishHandler(handler.NewGetIndexesHandler(dal, config))
	routes["/subsonic/rest/getartist"] = factory.PublishHandler(handler.NewGetArtistHandler(dal))
	routes["/subsonic/rest/getalbum"] = factory.PublishHandler(handler.NewGetAlbumHandler(dal))
	routes["/subsonic/rest/getsong"] = factory.PublishHandler(handler.NewGetSongHandler(dal))
	routes["/subsonic/rest/getrandomsongs"] = factory.PublishHandler(handler.NewGetRandomSongsHandler(dal))
	routes["/subsonic/rest/getmusicfolders"] = factory.PublishHandler(handler.NewGetMusicFoldersHandler(providers))
	routes["/subsonic/rest/getmusicdirectory"] = factory.PublishHandler(handler.NewGetMusicDirectoryHandler(dal))
	routes["/subsonic/rest/getgenres"] = factory.PublishHandler(handler.NewGetGenresHandler(dal))
	routes["/subsonic/rest/getsongsbygenre"] = factory.PublishHandler(handler.NewGetSongsByGenreHandler(dal))
	routes["/subsonic/rest/search2"] = factory.PublishHandler(handler.NewSearchHandler(dal, handler.Search2))
	routes["/subsonic/rest/search3"] = factory.PublishHandler(handler.NewSearchHandler(dal, handler.Search3))

	// Users
	routes["/subsonic/rest/getusers"] = factory.PublishHandler(handler.NewGetUsersHandler(config))
	routes["/subsonic/rest/getuser"] = factory.PublishHandler(handler.NewGetUserHandler(config))

	// Data
	routes["/subsonic/rest/getcoverart"] = factory.PublishBinaryHandler(handler.NewGetCoverArtHandler(dal))
	routes["/subsonic/rest/stream"] = factory.PublishBinaryHandler(handler.NewStreamHandler(dal))
	routes["/subsonic/rest/download"] = factory.PublishBinaryHandler(handler.NewDownloadHandler(dal))

	// Playlists
	routes["/subsonic/rest/createplaylist"] = factory.PublishHandler(handler.NewCreatePlaylistHandler(dal))
	routes["/subsonic/rest/getplaylists"] = factory.PublishHandler(handler.NewGetPlaylistsHandler(dal))
	routes["/subsonic/rest/getplaylist"] = factory.PublishHandler(handler.NewGetPlaylistHandler(dal))
	routes["/subsonic/rest/deleteplaylist"] = factory.PublishHandler(handler.NewDeletePlaylistHandler(dal))
	routes["/subsonic/rest/updateplaylist"] = factory.PublishHandler(handler.NewUpdatePlaylistHandler(dal))

	routes["/subsonic/rest/star"] = factory.PublishHandler(handler.NewStarHandler(dal, true))
	routes["/subsonic/rest/unstar"] = factory.PublishHandler(handler.NewStarHandler(dal, false))

	return routes
}
