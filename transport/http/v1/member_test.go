package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HackRVA/memberserver/datastore/in_memory"
	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/pkg/mqtt"
	"github.com/HackRVA/memberserver/services/member"
	resourcemanager "github.com/HackRVA/memberserver/transport/mqtt/v1"

	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"github.com/sirupsen/logrus"
)

var testMemberStore = in_memory.In_memory{
	Members: map[string]models.Member{
		"test@test.com": {
			ID:        "0",
			Name:      "testuser",
			Email:     "test@test.com",
			RFID:      "rfid1",
			Level:     0,
			Resources: []models.MemberResource{},
		},
		"test1@test.com": {
			ID:        "1",
			Name:      "testuser1",
			Email:     "test1@test.com",
			RFID:      "rfid2",
			Level:     0,
			Resources: []models.MemberResource{},
		},
		"test": {
			ID:        "",
			Name:      "",
			Email:     "test",
			RFID:      "",
			Level:     0,
			Resources: []models.MemberResource{},
		},
	},
	Tiers: []models.Tier{
		{
			ID:   0,
			Name: "fake-inactive",
		},
		{
			ID:   1,
			Name: "fake-active",
		},
		{
			ID:   2,
			Name: "fake-premium",
		},
	},
}

type slackNotifier struct{}

func (s slackNotifier) Send(msg string) {}

type paymentProvider struct{}

func (p paymentProvider) GetSubscription(subscriptionID string) (status string, lastPaymentAmount string, lastPaymentTime time.Time, err error) {
	return
}

func (p paymentProvider) GetSubscriber(subscriptionID string) (name string, email string, err error) {
	return
}

func TestGetMember(t *testing.T) {
	rm := resourcemanager.New(mqtt.New(), &in_memory.In_memory{}, slackNotifier{})
	server := &MemberServer{rm, member.New(&testMemberStore, rm, paymentProvider{}, logrus.New()), union.New()}

	// convert all members from the store to a json byte array
	jsonByte, _ := json.Marshal(in_memory.MemberMapToSlice(testMemberStore.Members))

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
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(jsonByte),
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newGetMembersRequest()
			response := httptest.NewRecorder()

			server.GetMembers(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestGetMemberByEmail(t *testing.T) {
	rm := resourcemanager.New(mqtt.New(), &in_memory.In_memory{}, slackNotifier{})
	server := &MemberServer{rm, member.New(&testMemberStore, rm, paymentProvider{}, logrus.New()), union.New()}

	// convert all members from the store to a json byte array
	jsonByte, _ := json.Marshal(testMemberStore.Members["test@test.com"])

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
			expectedHTTPStatus: http.StatusPreconditionFailed,
			expectedResponse:   "invalid email\n",
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

			server.GetByEmail(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestAssignRFID(t *testing.T) {
	rm := resourcemanager.New(mqtt.New(), &in_memory.In_memory{}, slackNotifier{})
	server := &MemberServer{rm, member.New(&testMemberStore, rm, paymentProvider{}, logrus.New()), union.New()}

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
			expectedResponse:   "{\"id\":\"0\",\"name\":\"testuser\",\"email\":\"test@test.com\",\"rfid\":\"rfid1\",\"memberLevel\":0,\"resources\":[],\"subscriptionID\":\"\"}",
		},
		{
			TestName:           "should return bad request if we don't send a valid email",
			Email:              "",
			RFID:               "newrfid",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusNotFound,
			expectedResponse:   "unable to assign rfid\n",
		},
		{
			TestName:           "should return bad request if we don't send a valid rfid string",
			Email:              "test@test.com",
			RFID:               "",
			Resources:          []models.MemberResource{},
			expectedHTTPStatus: http.StatusPreconditionFailed,
			expectedResponse:   "not a valid rfid\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			request := newAssignRFIDRequest(tt.Email, tt.RFID)
			response := httptest.NewRecorder()

			server.AssignRFID(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestGetTiers(t *testing.T) {
	rm := resourcemanager.New(mqtt.New(), &in_memory.In_memory{}, slackNotifier{})
	server := &MemberServer{rm, member.New(&testMemberStore, rm, paymentProvider{}, logrus.New()), union.New()}

	// convert all members from the store to a json byte array
	jsonByte, _ := json.Marshal(testMemberStore.Tiers)

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

			server.GetTiers(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestNewMember(t *testing.T) {
	rm := resourcemanager.New(mqtt.New(), &in_memory.In_memory{}, slackNotifier{})
	server := &MemberServer{rm, member.New(&testMemberStore, rm, paymentProvider{}, logrus.New()), union.New()}

	newMember := models.Member{
		ID:    "testID",
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

			server.AddNewMember(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedResponse)
		})
	}
}

func TestUpdateMemberSubscriptionID(t *testing.T) {
	rm := resourcemanager.New(mqtt.New(), &in_memory.In_memory{}, slackNotifier{})
	server := &MemberServer{rm, member.New(&testMemberStore, rm, paymentProvider{}, logrus.New()), union.New()}

	expectedResponse, _ := json.Marshal(models.EndpointSuccess{
		Ack: true,
	})

	tests := []struct {
		TestName           string
		Setup              func()
		Update             models.UpdateMemberRequest
		Expected           models.Member
		expectedHTTPStatus int
		expectedResponse   string
	}{
		{
			TestName: "should return a valid response for a valid email",
			Setup: func() {
				if _, err := server.MemberService.Add(models.Member{
					Name:           "testUser",
					Email:          "testUser@email.com",
					SubscriptionID: "unmodified",
				}); err != nil {
					logrus.Error(err)
				}
			},
			Expected: models.Member{
				Name:           "testUser",
				Email:          "testUser@email.com",
				SubscriptionID: "modified",
			},
			Update: models.UpdateMemberRequest{
				FullName:       "testUser",
				SubscriptionID: "modified",
			},
			expectedHTTPStatus: http.StatusOK,
			expectedResponse:   string(expectedResponse),
		},
		{
			TestName: "should respond not found if member doesn't exist",
			Setup:    func() {},
			Expected: models.Member{
				Name:           "doesn't exist",
				Email:          "doesntexist@email.com",
				SubscriptionID: "",
			},
			Update: models.UpdateMemberRequest{
				FullName:       "doesn't exist",
				SubscriptionID: "",
			},
			expectedHTTPStatus: http.StatusNotFound,
			expectedResponse:   "error getting member by email\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			tt.Setup()
			request := newUpdateMemberRequest(tt.Update, tt.Expected.Email)
			response := httptest.NewRecorder()

			server.UpdateMemberByEmail(response, request)

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

func newUpdateMemberRequest(update models.UpdateMemberRequest, email string) *http.Request {
	reqBody, _ := json.Marshal(update)
	req, _ := http.NewRequest(http.MethodPut, "/api/member/email/"+email, bytes.NewReader(reqBody))
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
