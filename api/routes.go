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
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// swagger:response VersionResponse
type versionResponse struct {
	// Commit Hash
	//
	// Example: "ffff"
	Commit string `json:"commit"`
}

// GitCommit is populated by a golang build arg
var GitCommit string

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
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: getUserResponse
	rr.HandleFunc("/user", api.getUser)
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
	rr.HandleFunc("/member", api.rbac(api.getMembers, []UserRole{admin}))
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
	rr.HandleFunc("/member/new", api.rbac(api.addNewMember, []UserRole{admin}))
	// swagger:route GET /api/member/self member getMemberByEmailRequest
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
	rr.HandleFunc("/member/self", api.getCurrentUserMemberInfo)
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
	rr.HandleFunc("/member/email/{email}", api.rbac(api.getMemberByEmail, []UserRole{admin}))
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
	rr.HandleFunc("/member/slack/nonmembers", api.rbac(api.getNonMembersOnSlack, []UserRole{admin}))
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
	rr.HandleFunc("/member/tier", api.rbac(api.getTiers, []UserRole{admin}))
	// swagger:route POST /api/payments/refresh payments getRefreshPayments
	//
	// Refresh payment information
	//
	// Submits a request to update member status information
	//   This will reach out to paypal and pull down the latest
	//   transaction information and then evaluate each member's
	//   membership status
	//
	//  This should happen automatically every day, but if we decide we
	//   want to manually update it.  This will give us the option to do so.
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
	//       200: getPaymentRefreshResponse
	rr.HandleFunc("/payments/refresh", api.rbac(api.refreshPayments, []UserRole{admin}))
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
	rr.HandleFunc("/member/assignRFID/self", api.assignRFIDSelf).Methods(http.MethodPost)
	// swagger:route POST /api/member/assignRFID member setRFIDRequest
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
	rr.HandleFunc("/member/assignRFID", api.rbac(api.assignRFID, []UserRole{admin})).Methods(http.MethodPost)
	// swagger:route GET /api/version version Version
	//
	// Version
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
	//       200: VersionResponse
	r.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		var version versionResponse

		version.Commit = GitCommit

		w.Header().Set("Content-Type", "application/json")

		j, _ := json.Marshal(version)
		w.Write(j)
	}).Methods(http.MethodGet)
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
	rr.HandleFunc("/auth/login", api.authenticate).Methods(http.MethodPost)
	// swagger:route POST /api/auth/logout auth logoutRequest
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
	r.HandleFunc("/api/auth/logout", api.logout)
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
	r.HandleFunc("/api/auth/register", api.signup)
	return rr
}
