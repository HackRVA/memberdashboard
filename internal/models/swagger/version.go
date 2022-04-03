package swagger

import "memberserver/internal/models"

// swagger:response versionResponse
type versionResponse struct {
	// in: body
	Body models.VersionResponse
}
