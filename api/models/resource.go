package models

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
