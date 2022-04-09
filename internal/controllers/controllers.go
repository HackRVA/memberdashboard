package controllers

import (
	"memberserver/internal/controllers/auth"
	"memberserver/internal/datastore"
	"memberserver/internal/services/config"
	"memberserver/internal/services/member"
	"memberserver/internal/services/resourcemanager"
	"memberserver/internal/services/resourcemanager/mqttserver"

	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

// API endpoints
type API struct {
	db             datastore.DataStore
	ResourceServer resourceAPI
	VersionServer  *VersionServer
	MemberServer   *MemberServer
	ReportsServer  *ReportsServer
	UserServer     *UserServer
	AuthStrategy   union.Union
	JWTKeeper      jwt.SecretsKeeper
}

type resourceAPI struct {
	db              datastore.DataStore
	config          config.Config
	resourcemanager resourcemanager.ResourceManager
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
		MemberServer:  &MemberServer{rm, member.NewMemberService(store, rm), auth.AuthStrategy},
		ReportsServer: &ReportsServer{store},
		UserServer:    &userServer,
		AuthStrategy:  auth.AuthStrategy,
		JWTKeeper:     auth.JWTSecretsKeeper,
	}
}
