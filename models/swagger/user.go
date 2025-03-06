package swagger

import (
	"github.com/HackRVA/memberserver/models"
)

// tokenResponseBody for json response of signin
// swagger:response loginResponse
type tokenResponseBody struct {
	// in: body
	Body models.TokenResponse
}

// swagger:parameters loginRequest
type loginRequest struct {
	// in: body
	Body models.Credentials
}

// swagger:response getUserResponse
type userResponseBody struct {
	// in: body
	Body models.UserResponse
}

// swagger:parameters registerUserRequest
type userRegisterRequest struct {
	// in: body
	Body models.Credentials
}
