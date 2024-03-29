package routes

import (
	_ "embed"
	"net/http"

	api "github.com/HackRVA/memberserver/pkg/membermgr/controllers"
	"github.com/HackRVA/memberserver/pkg/membermgr/controllers/auth"
	"github.com/HackRVA/memberserver/pkg/membermgr/middleware/rbac"
	"github.com/HackRVA/memberserver/pkg/membermgr/ui"

	"github.com/gorilla/mux"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

type Router struct {
	UnAuthedRouter *mux.Router
	authedRouter   *mux.Router
	api            api.API
	authStrategy   union.Union
}

// setupMiddleWare must run before other routes, so, we give it a separate function
func setupMiddleware(authedRouter *mux.Router, userServer UserHTTPHandler, auth *auth.AuthController) *mux.Router {
	authedRouter.Use(auth.AuthMiddleware)
	return authedRouter
}

func New(api api.API, auth *auth.AuthController) Router {
	unAuthedRouter := mux.NewRouter()
	authedRouter := unAuthedRouter.PathPrefix("/api/").Subrouter()
	setupMiddleware(authedRouter, api.UserServer, auth)

	router := Router{
		UnAuthedRouter: unAuthedRouter,
		authedRouter:   authedRouter,
		api:            api,
		authStrategy:   auth.AuthStrategy,
	}

	router.RegisterRoutes(auth)

	return router
}

func (r *Router) mountFS() {
	r.UnAuthedRouter.PathPrefix("/").Handler(spaRouter{ui.UI})
}

func (r *Router) RegisterRoutes(auth *auth.AuthController) *mux.Router {
	accessControl := rbac.New(r.authStrategy)
	r.setupUserRoutes(r.api.UserServer, auth)
	r.setupMemberRoutes(r.api.MemberServer, accessControl)
	r.setupResourceRoutes(r.api.ResourceServer, accessControl)
	r.setupPaymentRoutes(r.api, accessControl)
	r.setupReportsRoutes(r.api.ReportsServer, accessControl)
	r.setupVersionRoutes(r.api.VersionServer)

	r.mountFS()

	http.Handle("/", r.UnAuthedRouter)
	http.Handle("/api/", r.authedRouter)

	return r.authedRouter
}
