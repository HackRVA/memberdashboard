package v1

import (
	"net/http"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/datastore"
	"github.com/HackRVA/memberserver/integrations"
	"github.com/HackRVA/memberserver/pkg/paypal/listener"
	"github.com/HackRVA/memberserver/services"
	"github.com/HackRVA/memberserver/services/member"
	"github.com/HackRVA/memberserver/services/report"
	"github.com/HackRVA/memberserver/transport/http/middleware/auth"
	"github.com/HackRVA/memberserver/transport/http/middleware/rbac"
	"github.com/gorilla/mux"

	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

type API struct {
	db             datastore.DataStore
	unAuthedRouter *mux.Router
	authedRouter   *mux.Router
	rbac.AccessControl
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

// setupMiddleWare must run before other routes, so, we give it a separate function
func setupMiddleware(authedRouter *mux.Router, authController *auth.AuthController) *mux.Router {
	authedRouter.Use(authController.AuthMiddleware)
	return authedRouter
}

// Setup - setup us up the routes
func New(unAuthedRouter *mux.Router, store datastore.DataStore, authController *auth.AuthController, rm services.Resource, pp integrations.PaymentProvider, log services.Logger) API {
	c := config.Get()

	authedRouter := unAuthedRouter.PathPrefix("/api/").Subrouter()
	setupMiddleware(authedRouter, authController)
	userServer := NewUserServer(store, c)

	return API{
		db:             store,
		authedRouter:   authedRouter,
		unAuthedRouter: unAuthedRouter,
		AccessControl:  rbac.New(authController.AuthStrategy),
		resourceAPI: resourceAPI{
			db:              store,
			config:          c,
			resourcemanager: rm,
			logger:          log,
		},
		VersionServer:  &VersionServer{NewInMemoryVersionStore()},
		ReportsServer:  &ReportsServer{report.Report{Store: store}, log},
		MemberServer:   &MemberServer{rm, member.New(store, rm, pp), authController.AuthStrategy},
		UserServer:     &userServer,
		AuthStrategy:   authController.AuthStrategy,
		JWTKeeper:      authController.JWTSecretsKeeper,
		logger:         log,
		AuthController: authController,
	}
}

func (r API) SetupRoutes() {
	r.authedRouter.HandleFunc("/member", r.Restrict(r.GetMembers, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/new", r.Restrict(r.AddNewMember, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/self", r.GetCurrentUser)
	r.authedRouter.HandleFunc("/member/{id}/status", r.Restrict(r.CheckStatus, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/{id}", r.Restrict(r.UpdateMemberByID, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPut)
	r.authedRouter.HandleFunc("/member/email/{email}", r.Restrict(r.MemberEmail, []rbac.UserRole{rbac.Admin})).Methods(http.MethodGet, http.MethodPut)
	r.authedRouter.HandleFunc("/member/slack/nonmembers", r.Restrict(r.GetNonMembersOnSlack, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/tier", r.Restrict(r.GetTiers, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/assignRFID/self", r.AssignRFIDSelf).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/member/assignRFID", r.Restrict(r.AssignRFID, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/member/{id}/credit", r.Restrict(r.SetCredited, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPut)
	webhook := listener.New(true)
	r.unAuthedRouter.HandleFunc("/api/paypal/subscription/new", webhook.WebhooksHandler(r.PaypalSubscriptionWebHook))

	r.authedRouter.HandleFunc("/reports/membercounts", r.Restrict(r.GetMemberCountsCharts, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/reports/access", r.Restrict(r.GetAccessStatsChart, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/reports/churn", r.GetMemberChurn)

	r.authedRouter.HandleFunc("/resource", r.Restrict(r.Resource, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPut, http.MethodDelete, http.MethodGet)
	r.authedRouter.HandleFunc("/resource/status", r.Restrict(r.Status, []rbac.UserRole{rbac.Admin})).Methods(http.MethodGet)
	r.authedRouter.HandleFunc("/resource/register", r.Restrict(r.Register, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/resource/member/bulk", r.Restrict(r.AddMultipleMembersToResource, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/resource/deleteacls", r.Restrict(r.DeleteResourceACL, []rbac.UserRole{rbac.Admin})).Methods(http.MethodDelete)
	r.authedRouter.HandleFunc("/resource/updateacls", r.Restrict(r.UpdateResourceACL, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/resource/open", r.Restrict(r.Open, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/resource/member", r.Restrict(r.RemoveMember, []rbac.UserRole{rbac.Admin})).Methods(http.MethodDelete)
	r.authedRouter.HandleFunc("/user", r.GetUser)
	r.authedRouter.HandleFunc("/auth/login", r.Login).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/auth/logout", r.Logout).Methods(http.MethodDelete)
	r.unAuthedRouter.HandleFunc("/api/auth/register", r.RegisterUser)
	r.unAuthedRouter.HandleFunc("/api/version", r.Version)
}
