package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"memberserver/config"
	"memberserver/datastore/dbstore.go"
	"memberserver/resourcemanager"
	"memberserver/resourcemanager/mqttserver"
)

// API endpoints
type API struct {
	db            *dbstore.DatabaseStore
	resource      resourceAPI
	VersionServer *VersionServer
	MemberServer  *MemberServer
	UserServer    *UserServer
}

type resourceAPI struct {
	db              *dbstore.DatabaseStore
	config          config.Config
	resourcemanager resourcemanager.ResourceManager
}

// Setup - setup us up the routes
func Setup(db *dbstore.DatabaseStore) *mux.Router {
	c, _ := config.Load()

	userServer := NewUserServer(db, c)
	rm := resourcemanager.NewResourceManager(mqttserver.NewMQTTServer(), db)

	api := API{
		db: db,
		resource: resourceAPI{
			db:     db,
			config: c,
		},
		VersionServer: &VersionServer{NewInMemoryVersionStore()},
		MemberServer:  &MemberServer{db, rm},
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
