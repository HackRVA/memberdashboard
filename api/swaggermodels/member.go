package api

import (
	"memberserver/api/models"
	"memberserver/database"
)

// swagger:parameters getMemberByEmailRequest
type getMemberByEmailRequest struct {
	// in:path
	Email string `json:"email"`
}

// swagger:response getMembersResponse
type getMembersResponse struct {
	// in: body
	Body []database.Member
}

// swagger:response getMemberResponse
type getMemberResponse struct {
	// in: body
	Body database.Member
}

// swagger:response getTierResponse
type getTierResponse struct {
	// in: body
	Body []database.Tier
}

// swagger:response setRFIDResponse
type setRFIDResponse struct {
	// in: body
	Body database.AssignRFIDRequest
}

// swagger:parameters setRFIDRequest
type setRFIDRequest struct {
	// in: body
	Body database.AssignRFIDRequest
}

// swagger:response getPaymentRefreshResponse
type getPaymentRefreshResponse struct {
	Body models.EndpointSuccess
}

// swagger:parameters addNewMemberRequest
type addNewMemberRequest struct {
	// in: body
	Body models.NewMember
}
