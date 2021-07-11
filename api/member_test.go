package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"memberserver/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMember(t *testing.T) {
	server := &MemberServer{&testMemberStore}

	// convert all members from the store to a json byte array
	jsonByte, _ := json.Marshal(memberMapToSlice(testMemberStore.members))

	tests := []struct {
		TestName           string
		ID                 string
		Name               string
		Email              string
		RFID               string
		Level              uint8
		Resources          []models.MemberResource
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			TestName:           "should return all members",
			ID:                 "0",
			Name:               "testuser",
			Email:              "test@test.com",
			RFID:               "rfid1",
			Level:              0,
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(jsonByte),
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newGetMembersRequest()
			response := httptest.NewRecorder()

			server.GetMembersHandler(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestGetMemberByEmail(t *testing.T) {
	server := &MemberServer{&testMemberStore}

	// convert all members from the store to a json byte array
	jsonByte, _ := json.Marshal(testMemberStore.members["test@test.com"])

	tests := []struct {
		TestName           string
		ID                 string
		Name               string
		Email              string
		RFID               string
		Level              uint8
		Resources          []models.MemberResource
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			TestName:           "should return a valid response for a known email",
			Email:              "test@test.com",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(jsonByte),
		},
		{
			TestName:           "should show not found if email is empty",
			Email:              "",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusNotFound,
			expectedResponse:   "error getting member by email\n",
		},
		{
			TestName:           "should show not found if email is not in the store",
			Email:              "somthingElse@email.com",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusNotFound,
			expectedResponse:   "error getting member by email\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newGetMemberByEmailRequest(tt.Email)
			response := httptest.NewRecorder()

			server.GetByEmailHandler(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestAssignRFID(t *testing.T) {
	server := &MemberServer{&testMemberStore}

	tests := []struct {
		TestName           string
		Email              string
		RFID               string
		Resources          []models.MemberResource
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			TestName:           "should return a valid response for a valid request",
			Email:              "test@test.com",
			RFID:               "newrfid",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   "{\"id\":\"0\",\"name\":\"testuser\",\"email\":\"test@test.com\",\"rfid\":\"newrfid\",\"memberLevel\":0,\"resources\":[]}",
		},
		{
			TestName:           "should return bad request if we don't send a valid email",
			Email:              "",
			RFID:               "newrfid",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusBadRequest,
			expectedResponse:   "unable to assign rfid\n",
		},
		{
			TestName:           "should return bad request if we don't send a valid rfid string",
			Email:              "test@test.com",
			RFID:               "",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusBadRequest,
			expectedResponse:   "not a valid rfid\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newAssignRFIDRequest(tt.Email, tt.RFID)
			response := httptest.NewRecorder()

			server.AssignRFIDHandler(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestGetTiers(t *testing.T) {
	server := &MemberServer{&testMemberStore}

	// convert all members from the store to a json byte array
	jsonByte, _ := json.Marshal(testMemberStore.tiers)

	tests := []struct {
		TestName           string
		Resources          []models.MemberResource
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			TestName:           "should return all tiers",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(jsonByte),
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newGetTiersRequest()
			response := httptest.NewRecorder()

			server.GetTiersHandler(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestNewMember(t *testing.T) {
	server := &MemberServer{&testMemberStore}

	newMember := models.Member{
		Email: "test1@test.com",
		RFID:  "newrfid",
	}

	jsonByte, _ := json.Marshal(newMember)

	tests := []struct {
		TestName           string
		Member             models.Member
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			TestName:           "should return a valid response for a valid email",
			Member:             newMember,
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(jsonByte),
		},
		{
			TestName:           "should also update rfid info",
			Member:             newMember,
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(jsonByte),
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newNewMemberRequest(tt.Member)
			response := httptest.NewRecorder()

			server.AddNewMemberHandler(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func newAssignRFIDRequest(email, rfid string) *http.Request {
	assignReq := models.AssignRFIDRequest{
		RFID:  rfid,
		Email: email,
	}

	reqBody, _ := json.Marshal(assignReq)
	req, _ := http.NewRequest(http.MethodGet, "/api/members", bytes.NewReader(reqBody))
	return req
}

func newGetMembersRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/members", nil)
	return req
}

func newGetMemberByEmailRequest(email string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/member/email/%s", email), nil)
	return req
}

func newNewMemberRequest(member models.Member) *http.Request {
	reqBody, _ := json.Marshal(member)
	req, _ := http.NewRequest(http.MethodGet, "/api//member/new", bytes.NewReader(reqBody))
	return req
}

func newGetTiersRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/api/member/tier", nil)
	return req
}
