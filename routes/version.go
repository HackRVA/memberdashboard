package routes

import (
	"net/http"
)

type VersionHTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func (r Router) setupVersionRoutes(versionServer VersionHTTPHandler) {
	r.UnAuthedRouter.HandleFunc("/api/version", versionServer.ServeHTTP)
}
