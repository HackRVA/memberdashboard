package swagger

import "memberserver/internal/models"

// swagger:response endpointSuccessResponse
type endpointSuccessResponse struct {
	Body models.EndpointSuccess
}
