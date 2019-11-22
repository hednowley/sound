package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	// Subsonic API routes
	r.HandleFunc("/subsonic/", func(w http.ResponseWriter, r *http.Request) {
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
	r.HandleFunc("/api/authenticate", factory.NewBinaryHandler(controller.NewAuthenticateController(authenticator)))
	r.HandleFunc("/api/ticket", factory.NewHandler(controller.NewTicketController(ticketer)))
	r.HandleFunc("/api/stream", factory.NewBinaryHandler(controller.NewStreamController(dal)))
	r.HandleFunc("/api/art", factory.NewBinaryHandler(controller.NewArtController(dal)))
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.AddClient(ticketer, dal, w, r)
	})

	// Serves the front-end
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			// if we failed to get the absolute path respond wi th a 400 bad request and stop
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// prepend the path with the path to the static directory
		path = filepath.Join("static", path)

		// check whether a file exists at the given path
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			// file does not exist, serve index.html
			http.ServeFile(w, r, "static/index.html")
			return
		} else if err != nil {
			// if we got an error (that wasn't that the file doesn't exist) stating the file, return a 500 internal server error and stop
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.FileServer(http.Dir("static")).ServeHTTP(w, r)
	})

	http.Handle("/", r)

	// Websocket endpoints
	hub.SetHandler("getArtists", handlers.MakeGetArtistsHandler(dal))
	hub.SetHandler("getArtist", handlers.MakeGetArtistHandler(dal))
	hub.SetHandler("getAlbum", handlers.MakeGetAlbumHandler(dal))
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
