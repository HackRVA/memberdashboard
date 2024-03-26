package routes

import (
	"net/http"

	"github.com/HackRVA/memberserver/pkg/membermgr/middleware/rbac"
)

type ResourceHTTPHandler interface {
	Resource(w http.ResponseWriter, req *http.Request)
	AddMultipleMembersToResource(w http.ResponseWriter, req *http.Request)
	RemoveMember(w http.ResponseWriter, req *http.Request)
	Register(w http.ResponseWriter, req *http.Request)
	Status(w http.ResponseWriter, req *http.Request)
	UpdateResourceACL(w http.ResponseWriter, req *http.Request)
	Open(w http.ResponseWriter, req *http.Request)
	DeleteResourceACL(w http.ResponseWriter, req *http.Request)
}

func (r Router) setupResourceRoutes(resource ResourceHTTPHandler, accessControl rbac.RBAC) {
	r.authedRouter.HandleFunc("/resource", accessControl.Restrict(resource.Resource, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPut, http.MethodDelete, http.MethodGet)
	r.authedRouter.HandleFunc("/resource/status", accessControl.Restrict(resource.Status, []rbac.UserRole{rbac.Admin})).Methods(http.MethodGet)
	r.authedRouter.HandleFunc("/resource/register", accessControl.Restrict(resource.Register, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/resource/member/bulk", accessControl.Restrict(resource.AddMultipleMembersToResource, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/resource/deleteacls", accessControl.Restrict(resource.DeleteResourceACL, []rbac.UserRole{rbac.Admin})).Methods(http.MethodDelete)
	r.authedRouter.HandleFunc("/resource/updateacls", accessControl.Restrict(resource.UpdateResourceACL, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/resource/open", accessControl.Restrict(resource.Open, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/resource/member", accessControl.Restrict(resource.RemoveMember, []rbac.UserRole{rbac.Admin})).Methods(http.MethodDelete)
}
