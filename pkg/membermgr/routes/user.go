package routes

import (
	"net/http"
)

type UserHTTPHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
}

type AuthHTTPHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func (r Router) setupUserRoutes(userServer UserHTTPHandler, auth AuthHTTPHandler) {
	// swagger:route GET /api/user user user
	//
	// Returns the current logged in user.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//     - bearerAuth:
	//     - basicAuth:
	//
	//     Responses:
	//       200: getUserResponse
	r.authedRouter.HandleFunc("/user", userServer.GetUser)
	// swagger:route POST /api/auth/login auth loginRequest
	//
	// Login
	//
	// Login accepts some json with the `email` and `password`
	//   and returns some json that has the token string
	//
	//     Consumes:
	//     - application/json
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200: loginResponse
	r.authedRouter.HandleFunc("/auth/login", auth.Login).Methods(http.MethodPost)
	// swagger:route DELETE /api/auth/logout auth logoutRequest
	//
	// Logout
	//
	// Logout does not require a request body
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200:
	r.authedRouter.HandleFunc("/auth/logout", auth.Logout).Methods(http.MethodDelete)
	// swagger:route POST /api/auth/register auth registerUserRequest
	//
	// Register a new user
	//
	// Register a new user of the app
	//
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Responses:
	//       200: endpointSuccessResponse
	r.UnAuthedRouter.HandleFunc("/api/auth/register", auth.RegisterUser)
}
