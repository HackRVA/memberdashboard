package api

import (
	"encoding/json"
	"net/http"
)

// InfoResponse response of info request
// swagger:response infoResponse
type InfoResponse struct {
	// Info Message
	//
	// Example: "{ "message": "hello, world!"}"
	Message string `json:"message"`
}

// Info - a simple hello world
func (a API) Info(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(InfoResponse{
		Message: "hello, world!",
	})
	w.Write(j)
}
