package v1

import (
	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/datastore"
	"github.com/HackRVA/memberserver/integrations"
	"github.com/HackRVA/memberserver/services"
	"github.com/HackRVA/memberserver/services/member"
	"github.com/HackRVA/memberserver/services/report"
	"github.com/HackRVA/memberserver/transport/http/middleware/auth"

	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

type API struct {
	db datastore.DataStore
	resourceAPI
	*VersionServer
	*MemberServer
	*ReportsServer
	*UserServer
	*auth.AuthController
	AuthStrategy union.Union
	JWTKeeper    jwt.SecretsKeeper
	logger       Logger
}

// Setup - setup us up the routes
func New(store datastore.DataStore, auth *auth.AuthController, rm services.Resource, pp integrations.PaymentProvider, log services.Logger) API {
	c := config.Get()

	userServer := NewUserServer(store, c)

	return API{
		db: store,
		resourceAPI: resourceAPI{
			db:              store,
			config:          c,
			resourcemanager: rm,
			logger:          log,
		},
		VersionServer: &VersionServer{NewInMemoryVersionStore()},
		ReportsServer: &ReportsServer{report.Report{Store: store}, log},
		MemberServer:  &MemberServer{rm, member.New(store, rm, pp, log), auth.AuthStrategy},
		UserServer:    &userServer,
		AuthStrategy:  auth.AuthStrategy,
		JWTKeeper:     auth.JWTSecretsKeeper,
		logger:        log,
    AuthController: auth,
	}
}

// func (a *API) setupRoutes() {
// }
