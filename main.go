package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/hednowley/sound/api/api"
	"github.com/hednowley/sound/api/controller"
	"github.com/hednowley/sound/config"
	"github.com/hednowley/sound/dal"
	"github.com/hednowley/sound/database"
	"github.com/hednowley/sound/interfaces"
	"github.com/hednowley/sound/provider"
	"github.com/hednowley/sound/services"
	subsonic "github.com/hednowley/sound/subsonic/api"
	subsonicRoutes "github.com/hednowley/sound/subsonic/routes"
	"github.com/hednowley/sound/ws"
	"github.com/hednowley/sound/ws/handlers"
	"go.uber.org/fx"
)

// registerRoutes starts listening for HTTP requests.
func registerRoutes(
	factory *api.HandlerFactory,
	authenticator *services.Authenticator,
	ticketer *ws.Ticketer,
	dal interfaces.DAL,
	hub interfaces.Hub,
	scanner *provider.Scanner,
	routes subsonicRoutes.Routes) {

	// Subsonic API routes
	http.HandleFunc("/subsonic/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.ToLower(strings.Split(r.URL.Path, ".")[0])

		for p, h := range routes {
			if p == path {
				h(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	})

	// Serves the front-end
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Endpoints for websocket negotiation
	http.HandleFunc("/api/authenticate", factory.NewBinaryHandler(controller.NewAuthenticateController(authenticator)))
	http.HandleFunc("/api/ticket", factory.NewHandler(controller.NewTicketController(ticketer)))
	http.HandleFunc("/api/stream", factory.NewBinaryHandler(controller.NewStreamController(dal)))
	http.HandleFunc("/api/art", factory.NewBinaryHandler(controller.NewArtController(dal)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.AddClient(ticketer, dal, w, r)
	})

	// Websocket endpoints
	hub.SetHandler("getArtists", handlers.MakeGetArtistsHandler(dal))
	hub.SetHandler("getArtist", handlers.MakeGetArtistHandler(dal))
	hub.SetHandler("getAlbum", handlers.MakeGetAlbumHandler(dal))
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
			ws.NewHub,
			provider.NewScanner,
			services.NewAuthenticator,
			ws.NewTicketer,
			subsonic.NewHandlerFactory,
			api.NewHandlerFactory,
			subsonicRoutes.NewRoutes,
		),
		fx.Invoke(setUpLogger, registerRoutes, start),
	)

	app.Run()
}
