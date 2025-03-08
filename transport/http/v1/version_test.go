package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HackRVA/memberserver/models"
)

func NewStubVersionStore() *StubVersionStore {
	return &StubVersionStore{}
}

var mockVersion = models.VersionResponse{
	Major: "major",
	Build: GitCommit,
}

func (i *StubVersionStore) GetVersion() []byte {
	j, _ := json.Marshal(mockVersion)

	return j
}

type StubVersionStore struct{}

func TestVersion(t *testing.T) {
	server := &VersionServer{NewStubVersionStore()}

	expectedVersion := models.VersionResponse{
		Major: mockVersion.Major,
		Build: "test",
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
				mockVersion = expectedVersion
			},
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(expectedVersionJSON),
		},
		{
			name: "should fail if we didn't capture the commit hash",
			setup: func() {
				mockVersion.Major = ``
				mockVersion.Build = ``
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

			server.Version(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func newGetVersionRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/version", nil)
	return req
}
