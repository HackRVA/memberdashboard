package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ok(writer http.ResponseWriter, result interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(result)
	writer.Write(response)
}

// preconditionFailed -- validation error response message function
// should ONLY be used when the request does not meet the API requirements
func preconditionFailed(writer http.ResponseWriter, validateMessage string) {
	http.Error(writer, errors.New(validateMessage).Error(), http.StatusPreconditionFailed)
}

func notFound(writer http.ResponseWriter, errorMessage string) {
	http.Error(writer, errors.New(errorMessage).Error(), http.StatusNotFound)
}

func badRequest(writer http.ResponseWriter, errorMessage string) {
	http.Error(writer, errors.New(errorMessage).Error(), http.StatusBadRequest)
}

func internalServerError(writer http.ResponseWriter, errorMessage string) {
	http.Error(writer, errors.New(errorMessage).Error(), http.StatusInternalServerError)
}
