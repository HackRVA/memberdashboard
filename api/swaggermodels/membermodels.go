package api

import (
	"memberserver/api/models"
	"memberserver/database"
)

// swagger:response getMemberResponse
type memberResponseBody struct {
	// in: body
	Body []database.Member
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

// swagger:parameters getMemberByIDRequest
type getMemberByIDRequest struct {
	// in: query
	ID string
}
