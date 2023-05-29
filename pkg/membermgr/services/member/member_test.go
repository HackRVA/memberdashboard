package member_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/HackRVA/memberserver/pkg/membermgr/datastore/in_memory"
	"github.com/HackRVA/memberserver/pkg/membermgr/models"
	"github.com/HackRVA/memberserver/pkg/membermgr/services/logger"
	"github.com/HackRVA/memberserver/pkg/membermgr/services/member"
)

func TestPaymentBeforeOneMonthAgo(t *testing.T) {
	store := &in_memory.In_memory{}
	m := member.NewMemberService(store, models.Member{}, logger.New())

	tests := []struct {
		payment     models.Payment
		expected    bool
		description string
	}{
		{
			expected:    false,
			description: "now is not before a month ago",
			payment: models.Payment{
				Amount: "amount doesn't matter",
				Time:   time.Now(),
			},
		},
		{
			expected:    true,
			description: "32 days ago is before a month ago",
			payment: models.Payment{
				Amount: "amount doesn't matter",
				Time:   time.Now().Add((time.Hour * 24) * -32),
			},
		},
		{
			expected:    false,
			description: "29 days ago is not before a month ago",
			payment: models.Payment{
				Amount: "amount doesn't matter",
				Time:   time.Now().Add((time.Hour * 24) * -29),
			},
		},
		{
			expected:    false,
			description: "5 days ago is not before a month ago",
			payment: models.Payment{
				Amount: "amount doesn't matter",
				Time:   time.Now().Add((time.Hour * 24) * -5),
			},
		},
	}

	for _, tt := range tests {
		isOneMonthAgo := m.PaymentIsBeforeOneMonthAgo(tt.payment)
		if isOneMonthAgo != tt.expected {
			t.Errorf("expected: %t, received: %t", tt.expected, isOneMonthAgo)
		}
	}
}

func TestIsMemberActive(t *testing.T) {
	store := &in_memory.In_memory{}

	tests := []struct {
		expected           bool
		description        string
		currentMemberLevel models.MemberLevel
	}{
		{
			expected:           true,
			description:        "classic is active",
			currentMemberLevel: models.Classic,
		},
		{
			expected:           true,
			description:        "standard is active",
			currentMemberLevel: models.Standard,
		},
		{
			expected:           true,
			description:        "premium is active",
			currentMemberLevel: models.Premium,
		},
		{
			expected:           false,
			description:        "inactive is not active",
			currentMemberLevel: models.Inactive,
		},
	}

	for _, tt := range tests {
		m := member.NewMemberService(
			store, models.Member{
				Level: uint8(tt.currentMemberLevel),
			}, logger.New())
		if m.IsActive() != tt.expected {
			t.Errorf("expected: %t, received: %t", tt.expected, m.IsActive())
		}
	}
}

func TestMemberService_Add(t *testing.T) {
	mockStore := &in_memory.In_memory{}
	memberSvc := member.New(mockStore, nil, nil, nil)

	newMember := models.Member{
		Name:  "Test User",
		Email: "test@example.com",
	}

	addedMember, err := memberSvc.Add(newMember)
	assert.NoError(t, err)
	assert.NotEmpty(t, addedMember)
	assert.Equal(t, newMember.Name, addedMember.Name)
	assert.Equal(t, newMember.Email, addedMember.Email)
	assert.NotEmpty(t, addedMember.ID)
}
