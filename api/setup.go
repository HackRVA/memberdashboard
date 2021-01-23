package api

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"memberserver/config"
	"memberserver/database"
)

// API endpoints
type API struct {
	db       *database.Database
	config   config.Config
	resource resourceAPI
}

type resourceAPI struct {
	db     *database.Database
	config config.Config
}

// Setup - setup us up the routes
func Setup() *mux.Router {
	c, _ := config.Load()
	db, err := database.Setup()
	if err != nil {
		log.Fatal(fmt.Errorf("error setting up db: %s", err))
	}

	api := API{
		config: c,
		db:     db,
		resource: resourceAPI{
			db:     db,
			config: c,
		},
	}

	r := mux.NewRouter()
	restRouter := registerRoutes(r, api)
	serveSwaggerUI(r)

	spa := spaHandler{staticPath: "./ui/dist/", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	http.Handle("/", r)
	http.Handle("/api/", restRouter)

	return r
}
