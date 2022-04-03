package openapi

import (
	"net/http"

	"github.com/flowchartsman/swaggerui"
	"github.com/gorilla/mux"

	_ "embed"
)

//go:embed swagger.json
var spec []byte

func ServeSwaggerUI(r *mux.Router) {
	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", swaggerui.Handler(spec)))
}
