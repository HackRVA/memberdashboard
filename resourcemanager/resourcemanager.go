package resourcemanager

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"strings"

	log "github.com/sirupsen/logrus"

	"memberserver/database"
)

// Resource manager keeps the resources up to date by
//  pushing new updates and checking in on their health

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

const (
	// StatusGood - the resource is online and up to date
	StatusGood = iota
	// StatusOutOfDate - the resource does not have the most up to date information
	StatusOutOfDate
	// StatusOffline - the resource is not reachable
	StatusOffline
)

// UpdateResourceACL pulls a resource's accesslist from the DB and pushes it to the resource
func UpdateResourceACL(r database.Resource) error {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
		return err
	}
	// get acl for that resource
	accessList, err := db.GetResourceACL(r)

	if err != nil {
		return err
	}

	updateRequest := &ACLUpdateRequest{}
	updateRequest.ACL = accessList

	j, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}
	log.Debugf("access list: %s", j)

	// publish the update to mqtt broker
	Publish(r.Name+"/update", j)

	db.Release()

	return nil
}

// UpdateResources - publish an MQTT message to add a member to the actual device
func UpdateResources() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	resources := db.GetResources()

	for _, r := range resources {
		members, _ := db.GetResourceACLWithMemberInfo(r)
		for _, m := range members {
			b, _ := json.Marshal(&AddMemberRequest{
				ResourceAddress: r.Address,
				Command:         "adduser",
				UserName:        m.Name,
				RFID:            m.RFID,
				AccessType:      1,
				ValidUntil:      -86400,
			})
			Publish(r.Name, string(b))

			time.Sleep(2 * time.Second)
		}
	}

	db.Release()
}

func DeleteResourceACL() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}
	resources := db.GetResources()

	for _, r := range resources {
		b, _ := json.Marshal(&DeleteMemberRequest{
			ResourceAddress: r.Address,
			Command:         "deletusers", // not a type-o this is how the command is defined in the rfid reader
		})
		Publish(r.Name, string(b))
	}
	db.Release()
}

// CheckStatus will publish an mqtt command that requests for a specific device to verify that
//   the resource has the correct and up to date access list
//   It will do this by hashing the list retrieved from the DB and comparing it
//   with the hash that the resource reports
func CheckStatus(r database.Resource) {
	Publish(r.Name+"/cmd", "aclhash")
}

func hash(accessList []string) string {
	h := sha1.New()
	h.Write([]byte(strings.Join(accessList[:], "\n")))
	bs := h.Sum(nil)

	log.Debug(strings.Join(accessList[:], "\n"))
	log.Debugf("%x\n", bs)
	return fmt.Sprintf("%x", bs)
}
