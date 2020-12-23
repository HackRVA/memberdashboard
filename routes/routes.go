package routes

import (
	"fmt"
	"net/http"

	"log"

	"github.com/gorilla/mux"

	"github.com/dfirebaugh/memberserver/database"
)

// Setup - setup us up the routes
func Setup() {
	var err error
	r := mux.NewRouter()
	api := &API{}
	api.DB, err = database.Setup()

	if err != nil {
		log.Fatal(fmt.Errorf("error setting up db: %s", err))
	}

	r.HandleFunc("/", api.info)
	r.HandleFunc("/api/status", api.getStatuses)
	r.HandleFunc("/api/resource", api.getResources)
	r.HandleFunc("/api/tier", api.getTiers)
	r.HandleFunc("/api/member", api.getMembers)

	http.Handle("/", r)
}
