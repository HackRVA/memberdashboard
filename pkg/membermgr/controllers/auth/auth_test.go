package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HackRVA/memberserver/pkg/membermgr/datastore/in_memory"
	"github.com/HackRVA/memberserver/pkg/membermgr/models"
)

func TestRegisterUser(t *testing.T) {
	db := &in_memory.In_memory{
		Members: make(map[string]models.Member),
	}
	db.Members["test"] = models.Member{
		Name:  "test",
		Email: "test",
	}

	server := AuthController{
		store: db,
	}

	tests := []struct {
		TestName            string
		userName            string
		resources           []string
		creds               models.Credentials
		expectedHTTPStastub int
		expectedResponse    string
	}{
		{
			TestName: "should register a user",
			userName: "doesnt exist",
			creds: models.Credentials{
				Email:    "doesnt exist",
				Password: "password",
			},
			expectedHTTPStastub: http.StatusOK,
			expectedResponse:    "{\"ack\":true}",
		},
		{
			TestName: "should fail to register a user if they already exist",
			userName: "test",
			creds: models.Credentials{
				Email:    "test",
				Password: "password",
			},
			expectedHTTPStastub: http.StatusBadRequest,
			expectedResponse:    "error registering user\n",
		},
		{
			TestName: "should fail if password isn't provided",
			userName: "test",
			creds: models.Credentials{
				Email:    "doesn't exist 1",
				Password: "",
			},
			expectedHTTPStastub: http.StatusBadRequest,
			expectedResponse:    "password must be longer\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newRegisterUserRequest(tt.creds)
			response := httptest.NewRecorder()

			// authInfo := auth.NewDefaultUser(tt.userName, tt.userName, tt.resources, nil)
			server.RegisterUser(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStastub)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func newRegisterUserRequest(creds models.Credentials) *http.Request {
	reqBody, _ := json.Marshal(creds)
	req, _ := http.NewRequest(http.MethodGet, "/api/user", bytes.NewReader(reqBody))
	return req
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
