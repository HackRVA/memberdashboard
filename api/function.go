package api

import (
	"encoding/json"
	"net/http"
)

func ok(writer http.ResponseWriter, result interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(result)
	writer.Write(response)
}
