package api

import (
	"memberserver/api/models"
	"net/http"
)

// GetInfo - a simple hello world
func (a API) GetInfo(w http.ResponseWriter, req *http.Request) {
	ok(w, models.InfoResponse{
		Message: "hello, world!",
	})
}
