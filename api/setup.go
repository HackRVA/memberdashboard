package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"memberserver/config"
	"memberserver/database"
)

// API endpoints
type API struct {
	db            *database.Database
	resource      resourceAPI
	VersionServer *VersionServer
	MemberServer  *MemberServer
	UserServer    *UserServer
}

type resourceAPI struct {
	db     *database.Database
	config config.Config
}

// Setup - setup us up the routes
func Setup(db *database.Database) *mux.Router {
	c, _ := config.Load()

	userServer := NewUserServer(db, c)

	api := API{
		db: db,
		resource: resourceAPI{
			db:     db,
			config: c,
		},
		VersionServer: &VersionServer{NewInMemoryVersionStore()},
		MemberServer:  &MemberServer{db},
		UserServer:    &userServer,
	}

	r := mux.NewRouter()
	restRouter := registerRoutes(r, api)
	serveSwaggerUI(r)
	//set up go guardian here
	setupGoGuardian(c, userServer)

	spa := spaHandler{staticPath: "./ui/dist/", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	http.Handle("/", r)
	http.Handle("/api/", restRouter)

	return r
}
