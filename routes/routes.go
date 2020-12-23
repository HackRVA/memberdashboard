package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Setup - setup us up the routes
func Setup() {
	r := mux.NewRouter()
	api := &API{}

	r.HandleFunc("/", api.info)
	r.HandleFunc("/api/status", api.getStatuses)
	r.HandleFunc("/api/resource", api.getResources)
	r.HandleFunc("/api/tier", api.getTiers)
	r.HandleFunc("/api/member", api.getMembers)

	http.Handle("/", r)
}
