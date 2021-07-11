package api

import (
	"encoding/json"
	"memberserver/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersion(t *testing.T) {
	server := &VersionServer{NewInMemoryVersionStore()}

	expectedVersion := models.VersionResponse{
		Major:  "1",
		Minor:  "0",
		Hotfix: "0",
		Build:  "test",
	}
	expectedVersionJSON, _ := json.Marshal(expectedVersion)

	tests := []struct {
		name               string
		version            models.VersionResponse
		expectedHTTPStatus int
		expectedResponse   string
		setup              func()
	}{
		{
			name: "should respond with the test version",
			setup: func() {
				GitCommit = `test`
			},
			version: models.VersionResponse{
				Major:  "0",
				Minor:  "0",
				Hotfix: "0",
				Build:  "test",
			},
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(expectedVersionJSON),
		},
		{
			name: "should fail if we didn't capture the commit hash",
			setup: func() {
				GitCommit = ``
			},
			expectedHTTPStatus: http.StatusNotFound,
			expectedResponse:   "some issue getting the version",
		},
	}

	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			request := newGetVersionRequest()
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func newGetVersionRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/version", nil)
	return req
}
