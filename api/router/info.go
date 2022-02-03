package router

import (
	"net/http"
)

type InfoHTTPHandler interface {
	GetInfo(http.ResponseWriter, *http.Request)
}

func (r Router) setupInfoRoutes(infoServer InfoHTTPHandler) {
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
	r.authedRouter.HandleFunc("/info", infoServer.GetInfo)
}
