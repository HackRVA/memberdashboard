// Package Classification Member Server API.
//
//     Schemes: http, https
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
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

	"memberserver/payments/listener"

	"github.com/gorilla/mux"
)

func registerRoutes(r *mux.Router, api API) *mux.Router {
	rr := r.PathPrefix("/api/").Subrouter()
	rr.Use(api.UserServer.authMiddleware)
	// swagger:route GET /api/user user user
	//
	// Returns the current logged in user.
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
	//       200: getUserResponse
	rr.HandleFunc("/user", api.UserServer.getUser)
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
	rr.HandleFunc("/member", api.rbac(api.MemberServer.GetMembersHandler, []UserRole{admin}))
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
	rr.HandleFunc("/member/new", api.rbac(api.MemberServer.AddNewMemberHandler, []UserRole{admin}))
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
	rr.HandleFunc("/member/self", api.MemberServer.GetCurrentUserHandler)
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
	rr.HandleFunc("/member/email/{email}", api.rbac(api.MemberServer.MemberEmailHandler, []UserRole{admin})).Methods(http.MethodGet, http.MethodPut)
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
	rr.HandleFunc("/member/slack/nonmembers", api.rbac(api.MemberServer.GetNonMembersOnSlackHandler, []UserRole{admin}))
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
	rr.HandleFunc("/member/tier", api.rbac(api.MemberServer.GetTiersHandler, []UserRole{admin}))
	// swagger:route GET /api/payments/charts payments searchPaymentChartRequest
	//
	// Get Chart information of payments
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
	rr.HandleFunc("/payments/charts", api.rbac(api.getPaymentChart, []UserRole{admin}))
	// swagger:route GET /api/resource resource getResourceRequest
	//
	// Returns a resource.
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
	//       200: getResourceResponse

	// swagger:route PUT /api/resource resource updateResourceRequest
	//
	// Updates a resource.
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
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200:
	rr.HandleFunc("/resource", api.rbac(api.resource.Resource, []UserRole{admin})).Methods(http.MethodPut, http.MethodDelete, http.MethodGet)
	// swagger:route GET /api/resource/status resource getResourceStatus
	//
	// Returns status of the resources.
	//
	//  Returns the status of all resources.
	//    0 = Good
	//    1 = Out of Date
	//    2 = Offline
	//
	// if the resource is out of date, it will attempt to push an update
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
	//       200: getResourceStatusResponse
	rr.HandleFunc("/resource/status", api.rbac(api.resource.status, []UserRole{admin})).Methods(http.MethodGet)
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
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: postResourceResponse
	rr.HandleFunc("/resource/register", api.rbac(api.resource.register, []UserRole{admin})).Methods(http.MethodPost)
	// swagger:route POST /api/resource/member/bulk resource resourceBulkMemberRequest
	//
	// Adds multple members to a resource.
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
	//       200: addMulitpleMembersToResourceResponse
	rr.HandleFunc("/resource/member/bulk", api.rbac(api.resource.addMultipleMembersToResource, []UserRole{admin})).Methods(http.MethodPost)
	// swagger:route DELETE /api/resource/deleteacls resource resourceDeleteACLS
	//
	// Clears out all Resource ACLs on those devices
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
	rr.HandleFunc("/resource/deleteacls", api.rbac(api.resource.deleteResourceACL, []UserRole{admin})).Methods(http.MethodDelete)
	// swagger:route POST /api/resource/updateacls resource resourceUpdateACLS
	//
	// Attempts to send all members to a Resource
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
	rr.HandleFunc("/resource/updateacls", api.rbac(api.resource.updateResourceACL, []UserRole{admin})).Methods(http.MethodPost)
	// swagger:route POST /api/resource/open resource openResourceRequest
	//
	// sends an MQTT message to open a resource
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
	rr.HandleFunc("/resource/open", api.rbac(api.resource.open, []UserRole{admin})).Methods(http.MethodPost)
	// swagger:route DELETE /api/resource/member resource resourceRemoveMemberRequest
	//
	// Removes a member from a resource.
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
	//       200: removeMemberSuccessResponse
	rr.HandleFunc("/resource/member", api.rbac(api.resource.removeMember, []UserRole{admin})).Methods(http.MethodDelete)
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
	rr.HandleFunc("/info", api.Info)
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
	rr.HandleFunc("/member/assignRFID/self", api.MemberServer.AssignRFIDSelfHandler).Methods(http.MethodPost)
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
	rr.HandleFunc("/member/assignRFID", api.rbac(api.MemberServer.AssignRFIDHandler, []UserRole{admin})).Methods(http.MethodPost)
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
	r.HandleFunc("/api/version", api.VersionServer.ServeHTTP)
	// swagger:route POST /api/auth/login auth loginRequest
	//
	// Login
	//
	// Login accepts some json with the `email` and `password`
	//   and returns some json that has the token string
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
	//       200: loginResponse
	rr.HandleFunc("/auth/login", api.UserServer.login).Methods(http.MethodPost)
	// swagger:route DELETE /api/auth/logout auth logoutRequest
	//
	// Logout
	//
	// Logout does not require a request body
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200:
	r.HandleFunc("/api/auth/logout", api.UserServer.logout).Methods(http.MethodDelete)
	// swagger:route POST /api/auth/register auth registerUserRequest
	//
	// Register a new user
	//
	// Register a new user of the app
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200: endpointSuccessResponse
	r.HandleFunc("/api/auth/register", api.UserServer.registerUser)

	webhook := listener.New(true)
	r.HandleFunc("/api/paypal/subscription/new", webhook.WebhooksHandler(api.PaypalSubscriptionWebHookHandler))
	return rr
}
