// Package Classification Member Server API.
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Dustin Firebaugh<dafirebaugh@gmail.com> https://dustinfirebaugh.com
//
//    SecurityDefinitions:
//    bearerAuth:
//      type: apiKey
//      in: header
//      name: Authorization
//      description: Enter your bearer token
//
// swagger:meta
package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func registerRoutes(r *mux.Router, api API) *mux.Router {
	rr := r.PathPrefix("/api/").Subrouter()
	rr.Use(authMiddleware)
	// swagger:route GET /api/user user user
	//
	// Returns the current logged in user.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getUserResponse
	rr.HandleFunc("/user", api.getUser)
	// swagger:route GET /api/member member getMemberList
	//
	// Returns a list of the members in the system.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getMemberResponse
	rr.HandleFunc("/member", api.getMembers)
	// swagger:route GET /api/tier member getTiers
	//
	// Returns a list the member tiers.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getTierResponse
	rr.HandleFunc("/tier", api.getTiers)
	// swagger:route GET /api/resource resource getResourceRequest
	//
	// Returns a resource.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getResourceResponse

	// swagger:route POST /api/resource resource updateResourceRequest
	//
	// Updates a resource.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: postResourceResponse

	// swagger:route DELETE /api/resource resource deleteResourceRequest
	//
	// Deletes a resource.
	//
	//     Consumes:
	//     - application/json
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200:
	rr.HandleFunc("/resource", api.resource.Resource).Methods(http.MethodPost, http.MethodDelete, http.MethodGet)
	// swagger:route GET /api/resource/status resource getResourceStatus
	//
	// Returns status of the resources.
	//
	//  Returns the status of all resources.
	//    0 = Good
	//    1 = Out of Date
	//    3 = Offline
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getResourceStatusResponse
	rr.HandleFunc("/resource/status", api.resource.status).Methods(http.MethodGet)
	// swagger:route POST /api/resource/register resource registerResourceRequest
	//
	// Updates a resource.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: postResourceResponse
	rr.HandleFunc("/resource/register", api.resource.Resource).Methods(http.MethodPost)
	// swagger:route POST /api/resource/member/add resource resourceAddMemberRequest
	//
	// Adds a member to a resource.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: addMemberToResourceResponse
	rr.HandleFunc("/resource/member/add", api.resource.addMember).Methods(http.MethodPost)
	// swagger:route POST /api/resource/member/remove resource resourceRemoveMemberRequest
	//
	// Removes a member from a resource.
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: removeMemberToResourceResponse
	rr.HandleFunc("/resource/member/remove", api.resource.removeMember).Methods(http.MethodDelete)
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
	//     Schemes: http
	//
	//     Responses:
	//       200: infoResponse
	rr.HandleFunc("/info", api.Info)
	// swagger:route POST /api/login auth loginRequest
	//
	// Login
	//
	// Login accepts some json with the `username` and `password`
	//   and returns some json that has the token string
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Responses:
	//       200: loginResponse
	r.HandleFunc("/api/login", api.authenticate).Methods(http.MethodPost)
	// swagger:route POST /api/logout auth logoutRequest
	//
	// Logout
	//
	// Logout does not require a request body
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Responses:
	//       200:
	r.HandleFunc("/api/logout", api.logout)
	// swagger:route POST /api/register auth registerRequest
	//
	// Register a new user of the app
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Responses:
	//       200:
	r.HandleFunc("/api/register", api.signup)
	return rr
}
