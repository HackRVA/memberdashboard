package resourcemanager

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"memberserver/api/models"
	"memberserver/datastore"
	"memberserver/resourcemanager/mqttserver"
	"time"

	"strings"

	log "github.com/sirupsen/logrus"
)

var db datastore.DataStore

const (
	commandDeleteUID = "deletuid"
	commandAddUser   = "adduser"
	commandOpenDoor  = "opendoor"
)

// Resource manager keeps the resources up to date by
//  pushing new updates and checking in on their health

type ResourceManager struct {
	MQTTServer mqttserver.MQTTServer
	store      datastore.DataStore
}

const (
	// StatusGood - the resource is online and up to date
	StatusGood = iota
	// StatusOutOfDate - the resource does not have the most up to date information
	StatusOutOfDate
	// StatusOffline - the resource is not reachable
	StatusOffline
)

func NewResourceManager(ms mqttserver.MQTTServer, store datastore.DataStore) *ResourceManager {
	db = store
	return &ResourceManager{ms, store}
}

// UpdateResourceACL pulls a resource's accesslist from the DB and pushes it to the resource
func (rm *ResourceManager) UpdateResourceACL(r models.Resource) error {
	// get acl for that resource
	accessList, err := rm.store.GetResourceACL(r)

	if err != nil {
		return err
	}

	updateRequest := &models.ACLUpdateRequest{}
	updateRequest.ACL = accessList

	j, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}
	log.Debugf("access list: %s", j)

	// publish the update to mqtt broker
	rm.MQTTServer.Publish(r.Name+"/update", j)

	return nil
}

// UpdateResources - publish an MQTT message to add a member to the actual device
func (rm *ResourceManager) UpdateResources() {
	resources := rm.store.GetResources()

	for _, r := range resources {
		members, _ := rm.store.GetResourceACLWithMemberInfo(r)
		for _, m := range members {
			if m.Level == uint8(models.Inactive) {
				continue
			}

			b, _ := json.Marshal(&models.MemberRequest{
				ResourceAddress: r.Address,
				Command:         commandAddUser,
				UserName:        m.Name,
				RFID:            m.RFID,
				AccessType:      1,
				ValidUntil:      -86400,
			})
			rm.MQTTServer.Publish(r.Name, string(b))

			time.Sleep(2 * time.Second)
		}
	}
}

func (rm *ResourceManager) RemovedInvalidUIDs() {
	resources := rm.store.GetResources()

	log.Debug("looking for members to remove")

	for _, r := range resources {
		members := rm.store.GetMembers()
		for _, m := range members {
			if m.Level != uint8(models.Inactive) {
				return
			}

			if len(m.RFID) == 0 {
				return
			}

			/* We will just try to remove all invalid members even if they are already removed */
			b, _ := json.Marshal(&models.MemberRequest{
				ResourceAddress: r.Address,
				Command:         commandDeleteUID,
				RFID:            m.RFID,
			})
			rm.MQTTServer.Publish(r.Name, string(b))
			log.Debugf("attempting to remove member %s from rfid device %s", m.Email, r.Address)

			time.Sleep(2 * time.Second)
		}
	}
}

func (rm *ResourceManager) Open(resource models.Resource) {
	b, _ := json.Marshal(models.MQTTRequest{
		Door:    resource.Name,
		Command: commandOpenDoor,
	})

	rm.MQTTServer.Publish(resource.Name, string(b))
}

// PushOne - update one user on the resources
func (rm *ResourceManager) PushOne(m models.Member) {
	memberAccess, _ := rm.store.GetMembersAccess(m)
	for _, m := range memberAccess {
		b, _ := json.Marshal(&models.MemberRequest{
			ResourceAddress: m.ResourceAddress,
			Command:         commandAddUser,
			UserName:        m.Name,
			RFID:            m.RFID,
			AccessType:      1,
			ValidUntil:      -86400,
		})
		rm.MQTTServer.Publish(m.ResourceName, string(b))
	}
}

func (rm *ResourceManager) DeleteResourceACL() {
	resources := rm.store.GetResources()

	for _, r := range resources {
		b, _ := json.Marshal(&models.DeleteMemberRequest{
			ResourceAddress: r.Address,
			Command:         "deletusers", // not a type-o this is how the command is defined in the rfid reader
		})
		rm.MQTTServer.Publish(r.Name, string(b))
	}
}

// CheckStatus will publish an mqtt command that requests for a specific device to verify that
//   the resource has the correct and up to date access list
//   It will do this by hashing the list retrieved from the DB and comparing it
//   with the hash that the resource reports
func (rm *ResourceManager) CheckStatus(r models.Resource) {
	rm.MQTTServer.Publish(r.Name+"/cmd", "aclhash")
}

func hash(accessList []string) string {
	h := sha1.New()
	h.Write([]byte(strings.Join(accessList[:], "\n")))
	bs := h.Sum(nil)

	log.Debug(strings.Join(accessList[:], "\n"))
	log.Debugf("%x\n", bs)
	return fmt.Sprintf("%x", bs)
}
