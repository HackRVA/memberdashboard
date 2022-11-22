package swagger

import "github.com/HackRVA/memberserver/internal/models"

// swagger:response versionResponse
type versionResponse struct {
	// in: body
	Body models.VersionResponse
}
