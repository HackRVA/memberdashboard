package member_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/HackRVA/memberserver/datastore/in_memory"
	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/services/member"
)

type stubResourceManager struct{}

func (m *stubResourceManager) PushOne(member models.Member) {}

type MockPaymentProvider struct{}

func (m *MockPaymentProvider) GetSubscriber(subscriptionID string) (string, string, error) {
	if subscriptionID == "testSubID" {
		return "Test User", "test@example.com", nil
	}
	return "", "", fmt.Errorf("subscriber not found")
}

func (m *MockPaymentProvider) GetSubscription(subscriptionID string) (string, string, time.Time, error) {
	if subscriptionID == "testSubID" {
		return models.ActiveStatus, "10.00", time.Now(), nil
	}
	return "", "", time.Time{}, fmt.Errorf("subscription not found")
}

func TestMemberService_Add(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	memberService := member.New(mockStore, nil, nil)

	newMember := models.Member{
		Name:  "Test User",
		Email: "test@example.com",
	}

	addedMember, err := memberService.Add(newMember)
	assert.NoError(t, err)
	assert.NotEmpty(t, addedMember)
	assert.Equal(t, newMember.Name, addedMember.Name)
	assert.Equal(t, newMember.Email, addedMember.Email)
	assert.NotEmpty(t, addedMember.ID)
}

func TestMemberService_GetMembersPaginated(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	stubResourceManager := &stubResourceManager{}
	memberService := member.New(mockStore, stubResourceManager, nil)
	for i := 0; i < 10; i++ {
		if _, err := memberService.Add(models.Member{
			Name:           "Test User",
			Email:          "test" + strconv.Itoa(i) + "@example.com",
			RFID:           "abc1234" + strconv.Itoa(i),
			SubscriptionID: "abc1234" + strconv.Itoa(i),
			Level:          1,
		}); err != nil {
			t.Error(err)
		}
	}

	members := memberService.GetMembersPaginated(5, 0, true)
	assert.Len(t, members, 5)
}

func TestMemberService_GetByEmail(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	memberService := member.New(mockStore, nil, nil)

	newMember := models.Member{
		Name:  "Test User",
		Email: "test@example.com",
	}

	if _, err := memberService.Add(newMember); err != nil {
		t.Error(err)
	}
	member, err := memberService.GetByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, newMember.Email, member.Email)
}

func TestMemberService_Update(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	mockPaymentProvider := &MockPaymentProvider{}
	stubResourceManager := &stubResourceManager{}
	memberService := member.New(mockStore, stubResourceManager, mockPaymentProvider)

	newMember := models.Member{
		Name:           "Test User",
		Email:          "test@example.com",
		SubscriptionID: "testSubID",
		RFID:           "abc1234",
	}

	addedMember, _ := memberService.Add(newMember)
	addedMember.Name = "Updated User"
	err := memberService.Update(addedMember)
	assert.NoError(t, err)

	updatedMember, _ := memberService.GetByEmail("test@example.com")
	assert.Equal(t, "Updated User", updatedMember.Name)
}

func TestMemberService_AssignRFID(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	stubResourceManager := &stubResourceManager{}
	memberService := member.New(mockStore, stubResourceManager, nil)

	newMember := models.Member{
		Name:  "Test User",
		Email: "test@example.com",
		RFID:  "abc1234",
	}

	if _, err := memberService.Add(newMember); err != nil {
		t.Error(err)
	}
	assignedMember, err := memberService.AssignRFID("test@example.com", "123456")
	assert.NoError(t, err)
	assert.Equal(t, "123456", assignedMember.RFID)
}

func TestMemberService_GetMemberBySubscriptionID(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	mockPaymentProvider := &MockPaymentProvider{}
	memberService := member.New(mockStore, nil, mockPaymentProvider)

	newMember := models.Member{
		Name:           "Test User",
		Email:          "test@example.com",
		SubscriptionID: "testSubID",
	}

	if _, err := memberService.Add(newMember); err != nil {
		t.Error(err)
	}

	member, err := memberService.GetMemberBySubscriptionID("testSubID")
	assert.NoError(t, err)
	assert.Equal(t, newMember.Email, member.Email)
}

func TestMemberService_CheckStatus(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	mockPaymentProvider := &MockPaymentProvider{}
	memberService := member.New(mockStore, nil, mockPaymentProvider)

	newMember := models.Member{
		Name:           "Test User",
		Email:          "test@example.com",
		SubscriptionID: "testSubID",
		Level:          uint8(models.Standard),
	}

	if _, err := memberService.Add(newMember); err != nil {
		t.Error(err)
	}

	member, err := memberService.CheckStatus("testSubID")
	assert.NoError(t, err)
	assert.Equal(t, newMember.Email, member.Email)
}

func TestMemberService_FindNonMembersOnSlack(t *testing.T) {
	// This test requires a mock implementation of the Slack API
	// Assuming slack.GetUsers is mocked to return a list of users
}

func TestMemberService_SetLevel(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	memberService := member.New(mockStore, nil, nil)

	newMember := models.Member{
		Name:  "Test User",
		Email: "test@example.com",
	}

	addedMember, err := memberService.Add(newMember)
	assert.NoError(t, err)

	err = memberService.SetLevel(addedMember.ID, models.Premium)
	assert.NoError(t, err)

	updatedMember, _ := memberService.GetByEmail("test@example.com")
	assert.Equal(t, uint8(models.Premium), updatedMember.Level)
}

// func TestMemberService_GetActiveMembersWithoutSubscription(t *testing.T) {
// 	mockStore := &in_memory.In_memory{}
// 	memberService := member.New(mockStore, nil, nil)
//
// 	newMember := models.Member{
// 		Name:  "Test User",
// 		Email: "test@example.com",
// 		Level: uint8(models.Standard),
// 	}
//
// 	if _, err := memberService.Add(newMember); err != nil {
// 		t.Error(err)
// 	}
//
// 	members := memberService.GetActiveMembersWithoutSubscription()
// 	assert.Len(t, members, 1)
// 	assert.Equal(t, newMember.Email, members[0].Email)
// }
