package http

import (
	_ "embed"
	"io/fs"

	"github.com/HackRVA/memberserver/web"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type MemberServerRouter interface {
	SetupRoutes()
}

type Router struct {
	UnAuthedRouter *mux.Router
	routers        []MemberServerRouter
}

func New(unAuthedRouter *mux.Router, routers []MemberServerRouter) Router {
	router := Router{
		UnAuthedRouter: unAuthedRouter,
		routers:        routers,
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
	for _, router := range r.routers {
		router.SetupRoutes()
	}
	r.mountFS()
}
