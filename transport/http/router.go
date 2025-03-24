package http

import (
	"embed"
	"io/fs"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type MemberServerRouter interface {
	SetupRoutes()
}

type Router struct {
	UnAuthedRouter *mux.Router
	routers        []MemberServerRouter
	staticWebFS    embed.FS
}

func New(unAuthedRouter *mux.Router, routers []MemberServerRouter, staticWebFS embed.FS) Router {
	router := Router{
		UnAuthedRouter: unAuthedRouter,
		routers:        routers,
		staticWebFS:    staticWebFS,
	}

	router.registerRoutes()

	return router
}

func (r *Router) mountFS() {
	subFS, err := fs.Sub(r.staticWebFS, "web/dist/web-memberdashboard/browser")
	if err != nil {
		logrus.Fatal("Failed to mount static assets:", err)
	}

	r.UnAuthedRouter.PathPrefix("/").Handler(spaRouter{subFS})
}

func (r *Router) registerRoutes() {
	for _, router := range r.routers {
		router.SetupRoutes()
	}
	r.mountFS()
}
