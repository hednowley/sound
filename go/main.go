package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/controller"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/services"
	"github.com/hednowley/sound/socket"
	"github.com/hednowley/sound/socket/handlers"
	subsonic "github.com/hednowley/sound/subsonic/api"
	subsonicRoutes "github.com/hednowley/sound/subsonic/routes"
	"go.uber.org/fx"
)

// registerRoutes starts listening for HTTP requests.
func registerRoutes(
	factory *api.HandlerFactory,
	authenticator *services.Authenticator,
	ticketer *socket.Ticketer,
	dal *dal.DAL,
	hub socket.IHub,
	scanner *provider.Scanner,
	routes subsonicRoutes.Routes) {

	r := mux.NewRouter()

	// Subsonic API routes
	r.PathPrefix("/subsonic/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.ToLower(strings.Split(r.URL.Path, ".")[0])

		for p, h := range routes {
			if p == path {
				h(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	})

	// Endpoints for websocket negotiation
	r.HandleFunc("/api/authenticate", factory.NewHandler(controller.NewAuthenticateController(authenticator)))
	r.HandleFunc("/api/ticket", factory.NewHandler(controller.NewTicketController(ticketer)))
	r.HandleFunc("/api/stream", factory.NewBinaryHandler(controller.NewStreamController(dal)))
	r.HandleFunc("/api/art", factory.NewBinaryHandler(controller.NewArtController(dal)))
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.AddClient(ticketer, w, r)
	})

	// Serve the front-end
	r.PathPrefix("/").HandlerFunc(api.ServeSinglePageApp("static", "static/index.html"))

	http.Handle("/", r)

	// Websocket endpoints
	hub.SetHandler("getArtists", handlers.MakeGetArtistsHandler(dal))
	hub.SetHandler("getArtist", handlers.MakeGetArtistHandler(dal))
	hub.SetHandler("getAlbum", handlers.MakeGetAlbumHandler(dal))
	hub.SetHandler("getAlbums", handlers.MakeGetAlbumsHandler(dal))
	hub.SetHandler("getPlaylists", handlers.MakeGetPlaylistsHandler(dal))
	hub.SetHandler("getPlaylist", handlers.MakeGetPlaylistHandler(dal))
	hub.SetHandler("startScan", handlers.MakeStartScanHandler(scanner))

	go hub.Run()
}

func start(config *config.Config) {
	log.Error(http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil))
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
			config.NewConfig("config.yaml"),
			database.NewDefault,
			provider.NewProviders,
			dal.NewDAL,
			socket.NewHub,
			provider.NewScanner,
			services.NewAuthenticator,
			socket.NewTicketer,
			subsonic.NewHandlerFactory,
			api.NewHandlerFactory,
			subsonicRoutes.NewRoutes,
		),
		fx.Invoke(setUpLogger, registerRoutes, start),
	)

	app.Run()
}
