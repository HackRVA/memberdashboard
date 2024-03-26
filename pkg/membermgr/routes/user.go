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
	r.authedRouter.HandleFunc("/user", userServer.GetUser)
	r.authedRouter.HandleFunc("/auth/login", auth.Login).Methods(http.MethodPost)
	r.authedRouter.HandleFunc("/auth/logout", auth.Logout).Methods(http.MethodDelete)
	r.UnAuthedRouter.HandleFunc("/api/auth/register", auth.RegisterUser)
}
