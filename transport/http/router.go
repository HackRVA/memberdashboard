package http

import (
	_ "embed"
	"io/fs"
	"net/http"

	"github.com/HackRVA/memberserver/pkg/paypal/listener"
	"github.com/HackRVA/memberserver/transport/http/middleware/auth"
	"github.com/HackRVA/memberserver/transport/http/middleware/rbac"

	v1 "github.com/HackRVA/memberserver/transport/http/v1"

	"github.com/HackRVA/memberserver/web"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

type Router struct {
	UnAuthedRouter *mux.Router
	authedRouter   *mux.Router
	authStrategy   union.Union
	authController *auth.AuthController
	rbac.AccessControl

	v1.V1Router
}

// setupMiddleWare must run before other routes, so, we give it a separate function
func setupMiddleware(authedRouter *mux.Router, authController *auth.AuthController) *mux.Router {
	authedRouter.Use(authController.AuthMiddleware)
	return authedRouter
}

func New(v1Router v1.V1Router, authController *auth.AuthController) Router {
	unAuthedRouter := mux.NewRouter()
	authedRouter := unAuthedRouter.PathPrefix("/api/").Subrouter()
	setupMiddleware(authedRouter, authController)

	router := Router{
		UnAuthedRouter: unAuthedRouter,
		authedRouter:   authedRouter,
		authStrategy:   authController.AuthStrategy,
		V1Router:       v1Router,
		authController: authController,
		AccessControl:  rbac.New(authController.AuthStrategy),
	}

	router.registerRoutes()

	return router
}

func (r *Router) mountFS() {
	subFS, err := fs.Sub(web.UI, "dist/web-memberdashboard/browser")
	if err != nil {
		logrus.Fatal("Failed to mount static assets:", err)
	}

	r.UnAuthedRouter.PathPrefix("/").Handler(spaRouter{subFS})
}

func (r *Router) registerRoutes() {
	r.authedRouter.HandleFunc("/member", r.Restrict(r.GetMembers, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/new", r.Restrict(r.AddNewMember, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/self", r.GetCurrentUser)
	r.authedRouter.HandleFunc("/member/{id}/status", r.Restrict(r.CheckStatus, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/email/{email}", r.Restrict(r.MemberEmail, []rbac.UserRole{rbac.Admin})).Methods(http.MethodGet, http.MethodPut)
	r.authedRouter.HandleFunc("/member/slack/nonmembers", r.Restrict(r.GetNonMembersOnSlack, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/tier", r.Restrict(r.GetTiers, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/assignRFID/self", r.AssignRFIDSelf).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/member/assignRFID", r.Restrict(r.AssignRFID, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/member/{id}/credit", r.Restrict(r.SetCredited, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPut)
	webhook := listener.New(true)
	r.UnAuthedRouter.HandleFunc("/api/paypal/subscription/new", webhook.WebhooksHandler(r.PaypalSubscriptionWebHook))

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
	r.UnAuthedRouter.HandleFunc("/api/auth/register", r.RegisterUser)
	r.UnAuthedRouter.HandleFunc("/api/version", r.V1Router.ServeHTTP)

	r.mountFS()
}
