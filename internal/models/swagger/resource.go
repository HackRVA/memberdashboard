package swagger

import (
	"time"

	"github.com/HackRVA/memberserver/internal/models"
)

// swagger:parameters updateResourceRequest
type updateResourceRequest struct {
	// in: body
	Body models.ResourceRequest
}

// swagger:parameters registerResourceRequest
type registerResourceRequest struct {
	// in: body
	Body models.RegisterResourceRequest
}

// swagger:parameters deleteResourceRequest
type deleteResourceRequest struct {
	// in: body
	Body models.ResourceDeleteRequest
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

// swagger:parameters openResourceRequest
type openResourceRequest struct {
	// in: body
	Body models.OpenResourceRequest
}

// swagger:response getResourceResponse
type getResourceResponse struct {
	// in: body
	Body models.Resource
}

// swagger:response postResourceResponse
type postResourceResponse struct {
	// in: body
	Body models.Resource
}

// swagger:response addMemberToResourceResponse
type addMemberToResourceResponse struct {
	// in: body
	Body models.MemberResourceRelation
}

// swagger:response addMulitpleMembersToResourceResponse
type addMulitpleMembersToResourceResponse struct {
	// in: body
	Body []models.MemberResourceRelation
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
