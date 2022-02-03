package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type InfoHTTPHandler interface {
	GetInfo(http.ResponseWriter, *http.Request)
}

func setupInfoRoutes(unauthedRouter *mux.Router, authedRouter *mux.Router, infoServer InfoHTTPHandler) (*mux.Router, *mux.Router) {
	// swagger:route GET /api/resource resource getResourceRequest
	//

	// swagger:route GET /api/info info info
	//
	// A simple hello world.
	//
	// This will simply respond with some sample info
	//
	//     Produces:
	//     - application/json
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200: infoResponse
	authedRouter.HandleFunc("/info", infoServer.GetInfo)
	return unauthedRouter, authedRouter
}
