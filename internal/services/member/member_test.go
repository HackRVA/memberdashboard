package member_test

import (
	"memberserver/internal/datastore/in_memory"
	"memberserver/internal/models"
	"memberserver/internal/services/logger"
	"memberserver/internal/services/member"
	"testing"
	"time"
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
