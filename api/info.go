package api

import (
	"encoding/json"
	"memberserver/api/models"
	"net/http"
)

// Info - a simple hello world
func (a API) Info(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(models.InfoResponse{
		Message: "hello, world!",
	})
	w.Write(j)
}
