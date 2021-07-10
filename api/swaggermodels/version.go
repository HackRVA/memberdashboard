package api

import "memberserver/api/models"

// swagger:response versionResponse
type versionResponse struct {
	// in: body
	Body models.VersionResponse
}
