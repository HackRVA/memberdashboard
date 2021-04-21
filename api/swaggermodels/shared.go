package api

import "memberserver/api/models"

// swagger:response endpointSuccessResponse
type endpointSuccessResponse struct {
	Body models.EndpointSuccess
}
