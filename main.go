package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/hednowley/sound/api"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/handler"
	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/services"
	"go.uber.org/fx"
)

// registerHandlers associates routes with handlers.
func registerHandlers(factory *api.HandlerFactory, config *config.Config, db *dal.DAL, dal *dal.DAL) {

	handlers := make(map[string]http.HandlerFunc)

	handlers["/subsonic/rest/ping"] = factory.PublishHandler(handler.NewPingHandler())

	// Scanning
	handlers["/subsonic/rest/getscanstatus"] = factory.PublishHandler(handler.NewGetScanStatusHandler(dal))
	handlers["/subsonic/rest/startscan"] = factory.PublishHandler(handler.NewStartScanHandler(dal))
	handlers["/subsonic/rest/delete"] = factory.PublishHandler(handler.NewDeleteHandler(dal))

	// Querying
	handlers["/subsonic/rest/getalbumlist2"] = factory.PublishHandler(handler.NewGetAlbumList2Handler(db))
	handlers["/subsonic/rest/getartists"] = factory.PublishHandler(handler.NewGetArtistsHandler(db))
	handlers["/subsonic/rest/getindexes"] = factory.PublishHandler(handler.NewGetIndexesHandler(db))
	handlers["/subsonic/rest/getartist"] = factory.PublishHandler(handler.NewGetArtistHandler(db))
	handlers["/subsonic/rest/getalbum"] = factory.PublishHandler(handler.NewGetAlbumHandler(db))
	handlers["/subsonic/rest/getsong"] = factory.PublishHandler(handler.NewGetSongHandler(db))
	handlers["/subsonic/rest/getmusicdirectory"] = factory.PublishHandler(handler.NewGetMusicDirectoryHandler(db))
	handlers["/subsonic/rest/getgenres"] = factory.PublishHandler(handler.NewGetGenresHandler(db))
	handlers["/subsonic/rest/getsongsbygenre"] = factory.PublishHandler(handler.NewGetSongsByGenreHandler(db))

	// Users
	handlers["/subsonic/rest/getusers"] = factory.PublishHandler(handler.NewGetUsersHandler(config))
	handlers["/subsonic/rest/getuser"] = factory.PublishHandler(handler.NewGetUserHandler(config))

	// Data
	handlers["/subsonic/rest/getcoverart"] = factory.PublishBinaryHandler(handler.NewGetCoverArtHandler(db))
	handlers["/subsonic/rest/stream"] = factory.PublishBinaryHandler(handler.NewStreamHandler(db))
	handlers["/subsonic/rest/download"] = factory.PublishBinaryHandler(handler.NewDownloadHandler(db))

	// Playlists
	handlers["/subsonic/rest/createplaylist"] = factory.PublishHandler(handler.NewCreatePlaylistHandler(db))
	handlers["/subsonic/rest/getplaylists"] = factory.PublishHandler(handler.NewGetPlaylistsHandler(db))
	handlers["/subsonic/rest/getplaylist"] = factory.PublishHandler(handler.NewGetPlaylistHandler(db))
	handlers["/subsonic/rest/deleteplaylist"] = factory.PublishHandler(handler.NewDeletePlaylistHandler(db))
	handlers["/subsonic/rest/updateplaylist"] = factory.PublishHandler(handler.NewUpdatePlaylistHandler(db))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer log.Flush()
		log.Info(fmt.Sprintf("Request received: %v", r.URL.Path))
		path := strings.ToLower(strings.Split(r.URL.Path, ".")[0])

		for p, h := range handlers {
			if p == path {
				h(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	})

	defer log.Flush()
	log.Info(`********************`)
	log.Info("Application started!")
	log.Info(`********************`)

	log.Error(http.ListenAndServe(":"+config.Port, nil))
}

func setUpLogger(config *config.Config) {
	logger, err := log.LoggerFromConfigAsFile(config.LogConfig)
	if err != nil {
		fmt.Printf("ERROR: Could not initiate log. %v\n", err)
	} else {
		log.ReplaceLogger(logger)
	}
}

// Entry point for the application.
func main() {

	app := fx.New(
		fx.Provide(
			config.NewConfig,
			database.NewDefault,
			provider.NewProviders,
			dal.NewDAL,
			services.NewAuthenticator,
			api.NewHandlerFactory),
		fx.Invoke(setUpLogger, registerHandlers),
	)

	app.Run()
}
