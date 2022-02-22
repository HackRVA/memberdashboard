package router

import (
	"memberserver/api/rbac"
	"net/http"
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
	GetMemberCounts(http.ResponseWriter, *http.Request)
}

func (r Router) setupMemberRoutes(member MemberHTTPHandler, accessControl rbac.AccessControl) {
	// swagger:route GET /api/member member getMemberListRequest
	//
	// Returns a list of the members in the system.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getMembersResponse
	r.authedRouter.HandleFunc("/member", accessControl.Restrict(member.GetMembersHandler, []rbac.UserRole{rbac.Admin}))
	// swagger:route POST /api/member/new member addNewMemberRequest
	//
	// Add a new member
	//
	// Add a member that doesn't exist in our system.
	//  This would most likely be because they just signed up
	//  and we don't have information from paypal yet.
	//
	// If the paypal email doesn't match, their access will be revoked
	//   when we next sync with paypal.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: endpointSuccessResponse
	r.authedRouter.HandleFunc("/member/new", accessControl.Restrict(member.AddNewMemberHandler, []rbac.UserRole{rbac.Admin}))
	// swagger:route GET /api/member/self member getCurrentMemberRequest
	//
	// Returns current members information
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getMemberResponse
	r.authedRouter.HandleFunc("/member/self", member.GetCurrentUserHandler)
	// swagger:route GET /api/member/email/{email} member getMemberByEmailRequest
	//
	// Returns a member based on the email address.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getMemberResponse

	// swagger:route PUT /api/member/email/{email} member updateMemberRequest
	//
	// Updates a member.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: endpointSuccessResponse
	r.authedRouter.HandleFunc("/member/email/{email}", accessControl.Restrict(member.MemberEmailHandler, []rbac.UserRole{rbac.Admin})).Methods(http.MethodGet, http.MethodPut)
	// swagger:route GET /api/member/slack/nonmembers member getSlackNonMemberList
	//
	// Returns a list slack users that are possibly not members.
	//   It's entirely possible that these people are just using a different email than
	//   what they signed up with.  So, these accounts should be verified manually.
	//
	//     Produces:
	//     - test/csv
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: text/csv
	r.authedRouter.HandleFunc("/member/slack/nonmembers", accessControl.Restrict(member.GetNonMembersOnSlackHandler, []rbac.UserRole{rbac.Admin}))
	// swagger:route GET /api/member/tier member getTiers
	//
	// Returns a list the member tiers.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getTierResponse
	r.authedRouter.HandleFunc("/member/tier", accessControl.Restrict(member.GetTiersHandler, []rbac.UserRole{rbac.Admin}))
	// swagger:route POST /api/member/assignRFID/self member setRFIDRequest
	//
	// Assigns an RFID tag to the currently logged in user
	//
	//   it assigns an RFID tag to a member to the current user
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: setRFIDResponse
	r.authedRouter.HandleFunc("/member/assignRFID/self", member.AssignRFIDSelfHandler).Methods(http.MethodPost)
	// swagger:route POST /api/member/assignRFID member setSelfRFIDRequest
	//
	// Assigns an RFID tag to a member
	//
	//   it assigns an RFID tag to a member
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: setRFIDResponse
	r.authedRouter.HandleFunc("/member/assignRFID", accessControl.Restrict(member.AssignRFIDHandler, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
	// swagger:route GET /api/member/stats stats searchPaymentChartRequest
	//
	// Get Chart information of monthly member counts
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getPaymentChartResponse
	r.authedRouter.HandleFunc("/member/stats", accessControl.Restrict(member.GetMemberCounts, []rbac.UserRole{rbac.Admin}))
}
