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
	// swagger:route GET /api/member member getMemberList
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
	//       200: getMemberResponse
	rr.HandleFunc("/member", api.getMembers)
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
	rr.HandleFunc("/member/tier", api.getTiers)
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
	rr.HandleFunc("/payments/refresh", api.refreshPayments)
	// swagger:route GET /api/payments/charts payments getPaymentChart
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
	rr.HandleFunc("/payments/charts", api.getPaymentChart)
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
	rr.HandleFunc("/resource", api.resource.Resource).Methods(http.MethodPut, http.MethodDelete, http.MethodGet)
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
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//
	//     Responses:
	//       200: postResourceResponse
	rr.HandleFunc("/resource/register", api.resource.register).Methods(http.MethodPost)
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
	rr.HandleFunc("/resource/member/bulk", api.resource.addMultipleMembersToResource).Methods(http.MethodPost)
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
	rr.HandleFunc("/resource/deleteacls", api.resource.deleteResourceACL).Methods(http.MethodDelete)
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
	rr.HandleFunc("/resource/updateacls", api.resource.updateResourceACL).Methods(http.MethodPost)
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
	rr.HandleFunc("/resource/member", api.resource.removeMember).Methods(http.MethodDelete)
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
	// swagger:route POST /api/member/assignRFID member setRFIDRequest
	//
	// Assigns an RFID tag to a member
	//
	//   this is an unauthenticated request, for now.
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
	//     Responses:
	//       200: setRFIDResponse
	rr.HandleFunc("/member/assignRFID", api.assignRFID).Methods(http.MethodPost)
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
	r.HandleFunc("/api/auth/login", api.authenticate).Methods(http.MethodPost)
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
