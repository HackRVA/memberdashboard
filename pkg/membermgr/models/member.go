package models

import "time"

// NewMember - add a new member
type NewMember struct {
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
