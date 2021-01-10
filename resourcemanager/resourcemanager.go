package resourcemanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dfirebaugh/memberserver/database"
)

// Resource manager keeps the resources up to date by
//  pushing new updates and checking in on their health

// ResourceManager contains functions that
type ResourceManager struct {
	db *database.Database
}

// ACLUpdateRequest is the json object we send to a resource when pushing an update
type ACLUpdateRequest struct {
	ACL []string `json:"acl"`
}

// Setup initializes the resource manager
func Setup() *ResourceManager {
	var err error
	rm := &ResourceManager{}
	rm.db, err = database.Setup()

	if err != nil {
		log.Fatal(fmt.Errorf("error setting up db: %s", err))
	}

	return rm
}

// UpdateResourceACL pulls a resource's accesslist from the DB and pushes it to the resource
func (rm *ResourceManager) UpdateResourceACL(r database.Resource) error {
	// get acl for that resource
	accessList, err := rm.db.GetResourceACL(r)

	if err != nil {
		return err
	}

	updateRequest := &ACLUpdateRequest{}
	updateRequest.ACL = accessList

	j, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}

	// push the update to the resource
	resp, err := http.Post(r.Address+"/update", "application/json", bytes.NewBuffer(j))
	if err != nil {
		fmt.Println("Unable to reach the resource.")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)

		// TODO: check that the resource responds with a hash of the list
		fmt.Println("body=", string(body))
	}

	return nil
}
