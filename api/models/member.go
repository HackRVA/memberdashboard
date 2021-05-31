package models

// NewMember - add a new member
type NewMember struct {
	Email string `json:"email"`
	RFID  string `json:"rfid"`
}
