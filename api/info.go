package api

import (
	"memberserver/api/models"
	"net/http"
)

// Info - a simple hello world
func (a API) Info(w http.ResponseWriter, req *http.Request) {
	ok(w, models.InfoResponse{
		Message: "hello, world!",
	})
}
