package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/hednowley/sound/api/controller"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/services"
	"github.com/hednowley/sound/subsonic/api"
	"github.com/hednowley/sound/subsonic/handler"
	"github.com/hednowley/sound/ws"
	"github.com/hednowley/sound/ws/handlers"
	"go.uber.org/fx"

	api2 "github.com/hednowley/sound/api/api"
)

// registerSubsonicHandlers associates routes with handlers.
func registerSubsonicHandlers(factory *api.HandlerFactory, config *config.Config, dal interfaces.DAL, scanner *provider.Scanner) {

	handlers := make(map[string]http.HandlerFunc)

	handlers["/subsonic/rest/ping"] = factory.PublishHandler(handler.NewPingHandler())

	// Scanning
	handlers["/subsonic/rest/getscanstatus"] = factory.PublishHandler(handler.NewGetScanStatusHandler(scanner))
	handlers["/subsonic/rest/startscan"] = factory.PublishHandler(handler.NewStartScanHandler(scanner))
	handlers["/subsonic/rest/delete"] = factory.PublishHandler(handler.NewDeleteHandler(dal))

	// Querying
	handlers["/subsonic/rest/getalbumlist2"] = factory.PublishHandler(handler.NewGetAlbumList2Handler(dal))
	handlers["/subsonic/rest/getartists"] = factory.PublishHandler(handler.NewGetArtistsHandler(dal))
	handlers["/subsonic/rest/getindexes"] = factory.PublishHandler(handler.NewGetIndexesHandler(dal))
	handlers["/subsonic/rest/getartist"] = factory.PublishHandler(handler.NewGetArtistHandler(dal))
	handlers["/subsonic/rest/getalbum"] = factory.PublishHandler(handler.NewGetAlbumHandler(dal))
	handlers["/subsonic/rest/getsong"] = factory.PublishHandler(handler.NewGetSongHandler(dal))
	handlers["/subsonic/rest/getmusicdirectory"] = factory.PublishHandler(handler.NewGetMusicDirectoryHandler(dal))
	handlers["/subsonic/rest/getgenres"] = factory.PublishHandler(handler.NewGetGenresHandler(dal))
	handlers["/subsonic/rest/getsongsbygenre"] = factory.PublishHandler(handler.NewGetSongsByGenreHandler(dal))

	// Users
	handlers["/subsonic/rest/getusers"] = factory.PublishHandler(handler.NewGetUsersHandler(config))
	handlers["/subsonic/rest/getuser"] = factory.PublishHandler(handler.NewGetUserHandler(config))

	// Data
	handlers["/subsonic/rest/getcoverart"] = factory.PublishBinaryHandler(handler.NewGetCoverArtHandler(dal))
	handlers["/subsonic/rest/stream"] = factory.PublishBinaryHandler(handler.NewStreamHandler(dal))
	handlers["/subsonic/rest/download"] = factory.PublishBinaryHandler(handler.NewDownloadHandler(dal))

	// Playlists
	handlers["/subsonic/rest/createplaylist"] = factory.PublishHandler(handler.NewCreatePlaylistHandler(dal))
	handlers["/subsonic/rest/getplaylists"] = factory.PublishHandler(handler.NewGetPlaylistsHandler(dal))
	handlers["/subsonic/rest/getplaylist"] = factory.PublishHandler(handler.NewGetPlaylistHandler(dal))
	handlers["/subsonic/rest/deleteplaylist"] = factory.PublishHandler(handler.NewDeletePlaylistHandler(dal))
	handlers["/subsonic/rest/updateplaylist"] = factory.PublishHandler(handler.NewUpdatePlaylistHandler(dal))

	http.HandleFunc("/subsonic/", func(w http.ResponseWriter, r *http.Request) {
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
}

func registerAPIHandlers(factory *api2.HandlerFactory, config *config.Config, authenticator *services.Authenticator, ticketer *ws.Ticketer, dal interfaces.DAL, hub *ws.Hub, scanner *provider.Scanner) {

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/api/authenticate", factory.NewHandler(controller.NewAuthenticateController(authenticator, config)))
	http.HandleFunc("/api/ticket", factory.NewHandler(controller.NewTicketController(ticketer)))

	hub.SetHandler("getArtists", handlers.MakeGetArtistsHandler(dal))
	hub.SetHandler("startScan", handlers.MakeStartScanHandler(scanner))

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.NewClient(hub, ticketer, dal, w, r)
	})
}

func start(config *config.Config) {
	log.Error(http.ListenAndServe(":"+config.Port, nil))
	log.Info(`********************`)
	log.Info("Application started!")
	log.Info(`********************`)
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
			ws.NewHub,
			provider.NewScanner,
			services.NewAuthenticator,
			ws.NewTicketer,
			api.NewHandlerFactory,
			api2.NewHandlerFactory),
		fx.Invoke(setUpLogger,
			registerSubsonicHandlers,
			registerAPIHandlers, start),
	)

	app.Run()
}
