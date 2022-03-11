// Package Classification Member Server API.
//
//     Schemes: http, https
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//
//    SecurityDefinitions:
//    bearerAuth:
//      type: apiKey
//      in: header
//      name: Authorization
//      description: Enter your bearer token
//    basicAuth:
//      type: basic
//      in: header
//      name: Authorization
//      description: Enter your basic auth credentials
//
// swagger:meta
package routes

import (
	_ "embed"
	"memberserver/api/openapi"
	api "memberserver/internal/controllers"
	"memberserver/internal/controllers/auth"
	"memberserver/internal/middleware/rbac"
	"net/http"

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

func (r *Router) RegisterRoutes(auth *auth.AuthController) *mux.Router {
	accessControl := rbac.New(r.authStrategy)
	r.setupUserRoutes(r.api.UserServer, auth)
	r.setupMemberRoutes(r.api.MemberServer, accessControl)
	r.setupResourceRoutes(r.api.ResourceServer, accessControl)
	r.setupPaymentRoutes(r.api, accessControl)
	r.setupReportsRoutes(r.api.ReportsServer, accessControl)
	r.setupVersionRoutes(r.api.VersionServer)

	spa := spaHandler{staticPath: "./ui/dist/", indexPath: "index.html"}
	openapi.ServeSwaggerUI(r.UnAuthedRouter)
	r.UnAuthedRouter.PathPrefix("/").Handler(spa)
	http.Handle("/", r.UnAuthedRouter)
	http.Handle("/api/", r.authedRouter)

	return r.authedRouter
}
