package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func serveSwaggerUI(r *mux.Router) {
	fs := http.FileServer(http.Dir("./docs/swaggerui/"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))
}
