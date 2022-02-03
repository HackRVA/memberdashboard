package api

import (
	"memberserver/api/auth"
	"memberserver/config"
	"memberserver/datastore"
	"memberserver/resourcemanager"
	"memberserver/resourcemanager/mqttserver"

	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

// type Router interface {
// 	RegisterRoutes(r *mux.Router, api API, authStrategy union.Union) *mux.Router
// }

// API endpoints
type API struct {
	db             datastore.DataStore
	ResourceServer resourceAPI
	VersionServer  *VersionServer
	MemberServer   *MemberServer
	UserServer     *UserServer
	AuthStrategy   union.Union
	JWTKeeper      jwt.SecretsKeeper
}

type resourceAPI struct {
	db              datastore.DataStore
	config          config.Config
	resourcemanager *resourcemanager.ResourceManager
}

// Setup - setup us up the routes
func Setup(store datastore.DataStore, auth *auth.AuthController) API {
	c, _ := config.Load()

	userServer := NewUserServer(store, c)
	rm := resourcemanager.NewResourceManager(mqttserver.NewMQTTServer(), store)

	return API{
		db: store,
		ResourceServer: resourceAPI{
			db:              store,
			config:          c,
			resourcemanager: rm,
		},
		VersionServer: &VersionServer{NewInMemoryVersionStore()},
		MemberServer:  &MemberServer{store, rm, auth.AuthStrategy},
		UserServer:    &userServer,
		AuthStrategy:  auth.AuthStrategy,
		JWTKeeper:     auth.JWTSecretsKeeper,
	}
}
