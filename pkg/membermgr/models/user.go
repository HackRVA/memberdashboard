package models

import "github.com/golang-jwt/jwt"

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

// Credentials Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	// Password - the user's password
	// required: true
	// example: string
	Password string `json:"password"`
	// Email - the users email
	// required: true
	// example: string
	Email string `json:"email"`
}

// UserResponse - a user object that we can send as json
type UserResponse struct {
	// Email - user's Email
	// example: string
	Email     string     `json:"email"`
	Resources []Resource `json:"resources"`
}
