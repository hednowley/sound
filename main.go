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

	handlers["/rest/ping"] = factory.PublishHandler(handler.NewPingHandler())

	// Scanning
	handlers["/rest/getscanstatus"] = factory.PublishHandler(handler.NewGetScanStatusHandler(dal))
	handlers["/rest/startscan"] = factory.PublishHandler(handler.NewStartScanHandler(dal))
	handlers["/rest/delete"] = factory.PublishHandler(handler.NewDeleteHandler(dal))

	// Querying
	handlers["/rest/getalbumlist2"] = factory.PublishHandler(handler.NewGetAlbumList2Handler(db))
	handlers["/rest/getartists"] = factory.PublishHandler(handler.NewGetArtistsHandler(db))
	handlers["/rest/getindexes"] = factory.PublishHandler(handler.NewGetIndexesHandler(db))
	handlers["/rest/getartist"] = factory.PublishHandler(handler.NewGetArtistHandler(db))
	handlers["/rest/getalbum"] = factory.PublishHandler(handler.NewGetAlbumHandler(db))
	handlers["/rest/getsong"] = factory.PublishHandler(handler.NewGetSongHandler(db))
	handlers["/rest/getmusicdirectory"] = factory.PublishHandler(handler.NewGetMusicDirectoryHandler(db))
	handlers["/rest/getgenres"] = factory.PublishHandler(handler.NewGetGenresHandler(db))
	handlers["/rest/getsongsbygenre"] = factory.PublishHandler(handler.NewGetSongsByGenreHandler(db))

	// Users
	handlers["/rest/getusers"] = factory.PublishHandler(handler.NewGetUsersHandler(config))
	handlers["/rest/getuser"] = factory.PublishHandler(handler.NewGetUserHandler(config))

	// Data
	handlers["/rest/getcoverart"] = factory.PublishBinaryHandler(handler.NewGetCoverArtHandler(db))
	handlers["/rest/stream"] = factory.PublishBinaryHandler(handler.NewStreamHandler(db))
	handlers["/rest/download"] = factory.PublishBinaryHandler(handler.NewDownloadHandler(db))

	// Playlists
	handlers["/rest/createplaylist"] = factory.PublishHandler(handler.NewCreatePlaylistHandler(db))
	handlers["/rest/getplaylists"] = factory.PublishHandler(handler.NewGetPlaylistsHandler(db))
	handlers["/rest/getplaylist"] = factory.PublishHandler(handler.NewGetPlaylistHandler(db))
	handlers["/rest/deleteplaylist"] = factory.PublishHandler(handler.NewDeletePlaylistHandler(db))
	handlers["/rest/updateplaylist"] = factory.PublishHandler(handler.NewUpdatePlaylistHandler(db))

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
