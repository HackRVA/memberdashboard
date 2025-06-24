package models

import (
	"strings"
	"time"
)

// NewMember - add a new member
type NewMemberRequest struct {
	Email string `json:"email"`
	RFID  string `json:"rfid"`
}

// MemberResource a resource that a member belongs to
type MemberResource struct {
	ResourceID string `json:"resourceID"`
	Name       string `json:"name"`
}

// Member -- a member of the makerspace
type Member struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Email          string           `json:"email"`
	RFID           string           `json:"rfid"`
	Level          uint8            `json:"memberLevel"`
	Resources      []MemberResource `json:"resources"`
	SubscriptionID string           `json:"subscriptionID"`
}

// EnsureUpperSubscriptionID returns a copy of a Member with `strings.ToUpper`
// called on `SubscriptionID`.
func (m Member) EnsureUpperSubscriptionID() Member {
	mem := m
	mem.SubscriptionID = strings.ToUpper(m.SubscriptionID)
	return mem
}

// AssignRFIDRequest -- request to associate an rfid to a member
type AssignRFIDRequest struct {
	Email string `json:"email"`
	RFID  string `json:"rfid"`
}

// UpdateMemberRequest -- request to update a member
type UpdateMemberRequest struct {
	FullName       string `json:"fullName"`
	SubscriptionID string `json:"subscriptionID"`
}

type Payment struct {
	Amount string    `json:"amount"`
	Time   time.Time `json:"time"`
}

type MemberShipCreditRequest struct {
	IsCredited bool `json:"isCredited"`
}

type MembersPaginatedResponse struct {
	Members []Member `json:"members"`
	Count   uint `json:"count"`
}
