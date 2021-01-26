package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"memberserver/config"
	"memberserver/database"
	"memberserver/resourcemanager"
)

// API endpoints
type API struct {
	db       *database.Database
	config   config.Config
	resource resourceAPI
}

type resourceAPI struct {
	db     *database.Database
	rm     *resourcemanager.ResourceManager
	config config.Config
}

// Setup - setup us up the routes
func Setup(db *database.Database, rm *resourcemanager.ResourceManager) *mux.Router {
	c, _ := config.Load()

	api := API{
		config: c,
		db:     db,
		resource: resourceAPI{
			db:     db,
			rm:     rm,
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
