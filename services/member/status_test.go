package member_test

import (
	"testing"
	"time"

	"github.com/HackRVA/memberserver/datastore/in_memory"
	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/services/member"
)

func TestPaymentBeforeOneMonthAgo(t *testing.T) {
	store := &in_memory.In_memory{}
	m := models.Member{}
	statusChecker := member.NewStatusChecker(m, store, nil)

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
		isOneMonthAgo := statusChecker.PaymentIsBeforeOneMonthAgo(tt.payment)
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
		m := models.Member{
			Level: uint8(tt.currentMemberLevel),
		}
		statusChecker := member.NewStatusChecker(m, store, nil)
		if statusChecker.IsActive() != tt.expected {
			t.Errorf("expected: %t, received: %t", tt.expected, statusChecker.IsActive())
		}
	}
}

func TestIsMemberCredited(t *testing.T) {
	store := &in_memory.In_memory{}

	tests := []struct {
		expected           bool
		description        string
		currentMemberLevel models.MemberLevel
	}{
		{
			expected:           true,
			description:        "credited is credited",
			currentMemberLevel: models.Credited,
		},
		{
			expected:           false,
			description:        "classic is not credited",
			currentMemberLevel: models.Classic,
		},
		{
			expected:           false,
			description:        "standard is not credited",
			currentMemberLevel: models.Standard,
		},
		{
			expected:           false,
			description:        "premium is not credited",
			currentMemberLevel: models.Premium,
		},
	}

	for _, tt := range tests {
		m := models.Member{
			Level: uint8(tt.currentMemberLevel),
		}
		statusChecker := member.NewStatusChecker(m, store, nil)
		if statusChecker.IsCredited() != tt.expected {
			t.Errorf("expected: %t, received: %t", tt.expected, statusChecker.IsCredited())
		}
	}
}

func TestHasValidSubscriptionID(t *testing.T) {
	store := &in_memory.In_memory{}

	tests := []struct {
		expected       bool
		description    string
		subscriptionID string
	}{
		{
			expected:       true,
			description:    "valid subscription ID",
			subscriptionID: "validSubID",
		},
		{
			expected:       false,
			description:    "subscription ID is 'none'",
			subscriptionID: "none",
		},
		{
			expected:       false,
			description:    "empty subscription ID",
			subscriptionID: "",
		},
	}

	for _, tt := range tests {
		m := models.Member{
			SubscriptionID: tt.subscriptionID,
		}
		statusChecker := member.NewStatusChecker(m, store, nil)
		if statusChecker.HasValidSubscriptionID() != tt.expected {
			t.Errorf("expected: %t, received: %t", tt.expected, statusChecker.HasValidSubscriptionID())
		}
	}
}
