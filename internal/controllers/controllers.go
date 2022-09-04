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
	"memberserver/pkg/slack"

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
	resourcemanager resourcemanager.ResourceManager
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
func Setup(store datastore.DataStore, auth *auth.AuthController) API {
	c := config.Get()

	userServer := NewUserServer(store, c)
	log := logger.New()
	rm := resourcemanager.New(mqtt.New(), store, slack.Notifier{}, log)
	pp := paypal.Setup(c.PaypalURL, c.PaypalClientID, c.PaypalClientSecret, log)

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
