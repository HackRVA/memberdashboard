package controllers

import (
	"memberserver/internal/controllers/auth"
	"memberserver/internal/datastore"
	"memberserver/internal/services/config"
	"memberserver/internal/services/logger"
	"memberserver/internal/services/member"
	"memberserver/internal/services/report"
	"memberserver/internal/services/resourcemanager"
	"memberserver/pkg/mqtt"
	"memberserver/pkg/paypal"

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
	c := config.Get()

	userServer := NewUserServer(store, c)
	rm := resourcemanager.NewResourceManager(mqtt.New(), store)
	pp := paypal.Setup(c.PaypalURL, c.PaypalClientID, c.PaypalClientSecret, logger.New())

	return API{
		db: store,
		ResourceServer: resourceAPI{
			db:              store,
			config:          c,
			resourcemanager: rm,
		},
		VersionServer: &VersionServer{NewInMemoryVersionStore()},
		ReportsServer: &ReportsServer{report.Report{Store: store}},
		MemberServer:  &MemberServer{rm, member.NewMemberService(store, rm, pp), auth.AuthStrategy},
		UserServer:    &userServer,
		AuthStrategy:  auth.AuthStrategy,
		JWTKeeper:     auth.JWTSecretsKeeper,
	}
}
