package swagger

import "github.com/HackRVA/memberserver/pkg/membermgr/models"

// swagger:response versionResponse
type versionResponse struct {
	// in: body
	Body models.VersionResponse
}
