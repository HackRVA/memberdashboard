package controllers

import (
	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/pkg/membermgr/controllers/auth"
	"github.com/HackRVA/memberserver/pkg/membermgr/datastore"
	"github.com/HackRVA/memberserver/pkg/membermgr/integrations"
	"github.com/HackRVA/memberserver/pkg/membermgr/services"
	"github.com/HackRVA/memberserver/pkg/membermgr/services/member"
	"github.com/HackRVA/memberserver/pkg/membermgr/services/report"

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
	logger         Logger
}

type resourceAPI struct {
	db              datastore.DataStore
	config          config.Config
	resourcemanager services.Resource
	logger          Logger
}

type Logger interface {
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Print(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Trace(args ...interface{})
}

// Setup - setup us up the routes
func Setup(store datastore.DataStore, auth *auth.AuthController, rm services.Resource, pp integrations.PaymentProvider, log services.Logger) API {
	c := config.Get()

	userServer := NewUserServer(store, c)

	return API{
		db: store,
		ResourceServer: resourceAPI{
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
	}
}
