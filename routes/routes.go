package routes

import (
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/gorilla/mux"

	"github.com/dfirebaugh/memberserver/config"
	"github.com/dfirebaugh/memberserver/database"
)

// Setup - setup us up the routes
func Setup() *mux.Router {
	var err error

	c, _ := config.Load(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))
	r := mux.NewRouter()

	api := &API{}
	api.resource = resourceAPI{}
	api.config = c
	api.db, err = database.Setup()

	// give the resource routes access to the db
	api.resource.db = api.db

	if err != nil {
		log.Fatal(fmt.Errorf("error setting up db: %s", err))
	}

	restRouter := r.PathPrefix("/api/").Subrouter()
	restRouter.HandleFunc("/user", api.getUser)

	restRouter.HandleFunc("/info", api.info)
	restRouter.HandleFunc("/resource", api.resource.Resource).Methods(http.MethodPost, http.MethodDelete, http.MethodGet)
	restRouter.HandleFunc("/resource/member/add", api.resource.addMember).Methods(http.MethodPost)
	restRouter.HandleFunc("/resource/member/remove", api.resource.removeMember).Methods(http.MethodDelete)

	restRouter.HandleFunc("/tier", api.getTiers)
	restRouter.HandleFunc("/member", api.getMembers)
	restRouter.HandleFunc("/user", api.getUser)

	r.HandleFunc("/edge/register", api.Signup)
	r.HandleFunc("/edge/signin", api.Signin)
	r.HandleFunc("/edge/logout", api.authJWT(api.Logout))

	spa := spaHandler{staticPath: "./ui/dist/", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	http.Handle("/", r)

	return r
}
