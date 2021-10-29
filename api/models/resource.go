package models

import "time"

// MemberResourceRelationUpdateRequest -- update or delete a resource for a member
type MemberResourceRelationUpdateRequest struct {
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

// MemberAccess represents that a member has access to a certain resource.
//  this will get pushed to a device.
type MemberAccess struct {
	Email           string
	ResourceAddress string
	ResourceName    string
	Name            string
	RFID            string
}

// ACLUpdateRequest is the json object we send to a resource when pushing an update
type ACLUpdateRequest struct {
	ACL []string `json:"acl"`
}

// ACLResponse Response from a resource that is a hash of the ACL that the
//   resource has stored
type ACLResponse struct {
	Hash string `json:"acl"`
	// Name of the resource - this should match what we have in the database
	//  so we know which acl to compare it with
	Name string `json:"name"`
}

type MemberRequest struct {
	ResourceAddress string `json:"doorip"`
	Command         string `json:"cmd"`
	UserName        string `json:"user"`
	RFID            string `json:"uid"`
	AccessType      int    `json:"acctype"`
	ValidUntil      int    `json:"validuntil"`
}

type MQTTRequest struct {
	Command string `json:"cmd"`
}

type DeleteMemberRequest struct {
	ResourceAddress string `json:"doorip"`
	Command         string `json:"cmd"`
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

// ResourceDeleteRequest - request for deleting a resource
type ResourceDeleteRequest struct {
	// UniqueID of the Resource
	// required: true
	// example: string
	ID string `json:"id"`
}

// ResourceRequest a resource that can accept an access control list
type ResourceRequest struct {
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
	IsDefault bool `json:"isDefault"`
}

// RegisterResourceRequest a resource that can accept an access control list
type RegisterResourceRequest struct {
	// Name of the Resource
	// required: true
	// example: string
	Name string `json:"name"`
	// Address of the Resource. i.e. where it can be found on the network
	// required: true
	// example: string
	Address string `json:"address"`
	// Default state of the Resource
	// required: false
	// example: true
	IsDefault bool `json:"isDefault"`
}

// MemberResourceRelation  - a relationship between resources and members
type MemberResourceRelation struct {
	ID         string `json:"id"`
	MemberID   string `json:"memberID"`
	ResourceID string `json:"resourceID"`
}

// OpenResourceRequest -- request to associate an rfid to a member
type OpenResourceRequest struct {
	Name string `json:"name"`
}
