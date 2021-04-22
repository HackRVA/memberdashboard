package api

import (
	"memberserver/api/models"
	"memberserver/database"
	"time"
)

// swagger:parameters updateResourceRequest
type updateResourceRequest struct {
	// in: body
	Body database.ResourceRequest
}

// swagger:parameters registerResourceRequest
type registerResourceRequest struct {
	// in: body
	Body database.RegisterResourceRequest
}

// swagger:parameters deleteResourceRequest
type deleteResourceRequest struct {
	// in: body
	Body database.ResourceDeleteRequest
}

// swagger:parameters resourceAddMemberRequest
type resourceAddMemberRequest struct {
	// in: body
	Body models.MemberResourceRelation
}

// swagger:parameters resourceBulkMemberRequest
type resourceBulkMemberRequest struct {
	// in: body
	Body models.MembersResourceRelation
}

// swagger:response getResourceResponse
type getResourceResponse struct {
	// in: body
	Body database.Resource
}

// swagger:response postResourceResponse
type postResourceResponse struct {
	// in: body
	Body database.Resource
}

// swagger:response addMemberToResourceResponse
type addMemberToResourceResponse struct {
	// in: body
	Body database.MemberResourceRelation
}

// swagger:response addMulitpleMembersToResourceResponse
type addMulitpleMembersToResourceResponse struct {
	// in: body
	Body []database.MemberResourceRelation
}

// swagger:response getResourceStatusResponse
type getResourceStatusResponse struct {
	// in: body
	Body map[string]time.Time
}

// swagger:response removeMemberSuccessResponse
type removeMemberSuccessResponse struct {
	Body models.EndpointSuccess
}
