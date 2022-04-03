package controllers

import (
	"memberserver/internal/services/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shaj13/go-guardian/v2/auth"
)

func TestGetUser(t *testing.T) {
	c, _ := config.Load()
	server := NewUserServer(&testMemberStore, c)

	tests := []struct {
		TestName            string
		userName            string
		resources           []string
		expectedHTTPStastub int
		expectedResponse    string
	}{
		{
			TestName:            "should return currently logged in user",
			userName:            "test",
			expectedHTTPStastub: http.StatusOK,
			expectedResponse:    "{\"id\":\"\",\"name\":\"\",\"email\":\"test\",\"rfid\":\"\",\"memberLevel\":0,\"resources\":[],\"subscriptionID\":\"\"}",
		},
		{
			TestName:            "should return unauthorized if email doesn't exist",
			userName:            "doesnt exist",
			expectedHTTPStastub: http.StatusUnauthorized,
			expectedResponse:    "user not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newGetUserRequest()
			response := httptest.NewRecorder()

			authInfo := auth.NewDefaultUser(tt.userName, tt.userName, tt.resources, nil)
			server.GetUser(response, auth.RequestWithUser(authInfo, request))

			assertStatus(t, response.Code, tt.expectedHTTPStastub)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func newGetUserRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
	return req
}
