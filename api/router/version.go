package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type VersionHTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func setupVersionRoutes(unauthedRouter *mux.Router, authedRouter *mux.Router, versionServer VersionHTTPHandler) (*mux.Router, *mux.Router) {
	// swagger:route GET /api/version version Version
	//
	//   Shows the current build's version information
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200: versionResponse
	unauthedRouter.HandleFunc("/api/version", versionServer.ServeHTTP)
	return unauthedRouter, authedRouter
}
