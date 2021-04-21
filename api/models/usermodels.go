package models

import (
	"github.com/dgrijalva/jwt-go"
)

// Claims -- auth claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// TokenResponse -- for json response of signin
type TokenResponse struct {
	// login response to send token string
	//
	// Example: <TOKEN_STRING>
	Token string `json:"token"`
}
