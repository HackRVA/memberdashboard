package routes

import (
	"net/http"

	"github.com/HackRVA/memberserver/pkg/membermgr/middleware/rbac"
)

type MemberHTTPHandler interface {
	MemberEmailHandler(w http.ResponseWriter, r *http.Request)
	GetMembersHandler(w http.ResponseWriter, r *http.Request)
	UpdateMemberByEmailHandler(w http.ResponseWriter, r *http.Request)
	GetByEmailHandler(w http.ResponseWriter, r *http.Request)
	GetCurrentUserHandler(w http.ResponseWriter, r *http.Request)
	AssignRFIDHandler(w http.ResponseWriter, r *http.Request)
	AssignRFIDSelfHandler(w http.ResponseWriter, r *http.Request)
	GetTiersHandler(w http.ResponseWriter, r *http.Request)
	GetNonMembersOnSlackHandler(w http.ResponseWriter, r *http.Request)
	AddNewMemberHandler(w http.ResponseWriter, r *http.Request)
	CheckStatus(w http.ResponseWriter, r *http.Request)
	SetCredited(w http.ResponseWriter, r *http.Request)
}

func (r Router) setupMemberRoutes(member MemberHTTPHandler, accessControl rbac.AccessControl) {
	r.authedRouter.HandleFunc("/member", accessControl.Restrict(member.GetMembersHandler, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/new", accessControl.Restrict(member.AddNewMemberHandler, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/self", member.GetCurrentUserHandler)
	r.authedRouter.HandleFunc("/member/{id}/status", accessControl.Restrict(member.CheckStatus, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/email/{email}", accessControl.Restrict(member.MemberEmailHandler, []rbac.UserRole{rbac.Admin})).Methods(http.MethodGet, http.MethodPut)
	r.authedRouter.HandleFunc("/member/slack/nonmembers", accessControl.Restrict(member.GetNonMembersOnSlackHandler, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/tier", accessControl.Restrict(member.GetTiersHandler, []rbac.UserRole{rbac.Admin}))
	r.authedRouter.HandleFunc("/member/assignRFID/self", member.AssignRFIDSelfHandler).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/member/assignRFID", accessControl.Restrict(member.AssignRFIDHandler, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/member/{id}/credit", accessControl.Restrict(member.SetCredited, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPut)
}
