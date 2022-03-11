package routes

import (
	"net/http"
)

type VersionHTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func (r Router) setupVersionRoutes(versionServer VersionHTTPHandler) {
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
	r.UnAuthedRouter.HandleFunc("/api/version", versionServer.ServeHTTP)
}
