package models

import "time"

// MemberResourceRelation -- update or delete a resource for a member
type MemberResourceRelation struct {
	// ID of the Resource
	// required: true
	// example: string
	ID string `json:"resourceID"`
	// Email - this will be the member's email address
	// Name of the Resource
	// required: true
	// example: email
	Email string `json:"email"`
}

// MembersResourceRelation -- update or delete a resource for multiple members
type MembersResourceRelation struct {
	// ID of the Resource
	// required: true
	// example: string
	ID string `json:"resourceID"`
	// Emails - list of member's email address
	// required: true
	// example: []
	Emails []string `json:"emails"`
}

// Resource a resource that can accept an access control list
type Resource struct {
	// UniqueID of the Resource
	// required: true
	// example: string
	ID string `json:"id"`
	// Name of the Resource
	// required: true
	// example: string
	Name string `json:"name"`
	// Address of the Resource. i.e. where it can be found on the network
	// required: true
	// example: string
	Address string `json:"address"`
	// Default state of the Resource
	// required: true
	// example: true
	IsDefault     bool      `json:"isDefault"`
	LastHeartBeat time.Time `json:"lastHeartBeat"`
}
