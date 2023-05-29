package swagger

import "github.com/HackRVA/memberserver/pkg/membermgr/models"

// swagger:response endpointSuccessResponse
type endpointSuccessResponse struct {
	Body models.EndpointSuccess
}
