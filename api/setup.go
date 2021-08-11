package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"memberserver/config"
	"memberserver/datastore"
	"memberserver/resourcemanager"
	"memberserver/resourcemanager/mqttserver"
)

// API endpoints
type API struct {
	db            datastore.DataStore
	resource      resourceAPI
	VersionServer *VersionServer
	MemberServer  *MemberServer
	UserServer    *UserServer
}

type resourceAPI struct {
	db              datastore.DataStore
	config          config.Config
	resourcemanager resourcemanager.ResourceManager
}

// Setup - setup us up the routes
func Setup(store datastore.DataStore) *mux.Router {
	c, _ := config.Load()

	userServer := NewUserServer(store, c)
	rm := resourcemanager.NewResourceManager(mqttserver.NewMQTTServer(), store)

	api := API{
		db: store,
		resource: resourceAPI{
			db:     store,
			config: c,
		},
		VersionServer: &VersionServer{NewInMemoryVersionStore()},
		MemberServer:  &MemberServer{store, rm},
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
