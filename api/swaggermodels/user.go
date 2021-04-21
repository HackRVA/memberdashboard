package api

import (
	"memberserver/api/models"
	"memberserver/database"
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
	Body database.Credentials
}

// swagger:response getUserResponse
type userResponseBody struct {
	// in: body
	Body database.UserResponse
}

// swagger:parameters registerUserRequest
type userRegisterRequest struct {
	// in: body
	Body database.Credentials
}
