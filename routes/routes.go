package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Setup - setup us up the routes
func Setup() {
	r := mux.NewRouter()

	r.HandleFunc("/", info)
	r.HandleFunc("/api/status", getStatuses)
	r.HandleFunc("/api/resource", getResources)
	r.HandleFunc("/api/tier", getTiers)
	r.HandleFunc("/api/member", getMembers)

	http.Handle("/", r)
}
