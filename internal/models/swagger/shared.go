package swagger

import "github.com/HackRVA/memberserver/internal/models"

// swagger:response endpointSuccessResponse
type endpointSuccessResponse struct {
	Body models.EndpointSuccess
}
