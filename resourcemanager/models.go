package resourcemanager

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

type AddMemberRequest struct {
	ResourceAddress string `json:"doorip"`
	Command         string `json:"cmd"`
	UserName        string `json:"user"`
	RFID            string `json:"uid"`
	AccessType      int    `json:"acctype"`
	ValidUntil      int    `json:"validuntil"`
}

type DeleteMemberRequest struct {
	ResourceAddress string `json:"doorip"`
	Command         string `json:"cmd"`
}
