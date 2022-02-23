package router

import (
	"memberserver/api/rbac"
	"net/http"
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
	// swagger:route GET /api/resource resource updateResourceRequest
	// Returns a resource.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//     - basicAuth:
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
	//     - basicAuth:
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
	//     - basicAuth:
	//
	//     Responses:
	//       200:
	r.authedRouter.HandleFunc("/resource", accessControl.Restrict(resource.Resource, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPut, http.MethodDelete, http.MethodGet)
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
	//     - basicAuth:
	//
	//     Responses:
	//       200: getResourceStatusResponse
	r.authedRouter.HandleFunc("/resource/status", accessControl.Restrict(resource.Status, []rbac.UserRole{rbac.Admin})).Methods(http.MethodGet)
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
	//     - basicAuth:
	//
	//     Responses:
	//       200: postResourceResponse
	r.authedRouter.HandleFunc("/resource/register", accessControl.Restrict(resource.Register, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
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
	//     - basicAuth:
	//
	//     Responses:
	//       200: addMulitpleMembersToResourceResponse
	r.authedRouter.HandleFunc("/resource/member/bulk", accessControl.Restrict(resource.AddMultipleMembersToResource, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
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
	//     - basicAuth:
	//
	//     Responses:
	//       200: endpointSuccessResponse
	r.authedRouter.HandleFunc("/resource/deleteacls", accessControl.Restrict(resource.DeleteResourceACL, []rbac.UserRole{rbac.Admin})).Methods(http.MethodDelete)
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
	r.authedRouter.HandleFunc("/resource/updateacls", accessControl.Restrict(resource.UpdateResourceACL, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
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
	//     - basicAuth:
	//
	//     Responses:
	//       200: endpointSuccessResponse
	r.authedRouter.HandleFunc("/resource/open", accessControl.Restrict(resource.Open, []rbac.UserRole{rbac.Admin})).Methods(http.MethodPost)
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
	//     - basicAuth:
	//
	//     Responses:
	//       200: removeMemberSuccessResponse
	r.authedRouter.HandleFunc("/resource/member", accessControl.Restrict(resource.RemoveMember, []rbac.UserRole{rbac.Admin})).Methods(http.MethodDelete)
}
