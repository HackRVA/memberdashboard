package resourcemanager

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"

	"memberserver/api/models"
	"memberserver/database"
)

// Resource manager keeps the resources up to date by
//  pushing new updates and checking in on their health

type MQTTServer interface {
	Publish(topic string, payload interface{})
	Subscribe(topic string, handler mqtt.MessageHandler)
}

// RMStore is where the resourcemanager gets the data from.
//  in prod this is the database
//  but we can make something up for testing
type RMStore interface {
	GetResources() []models.Resource
	//GetResourceACL returns a csv of rfid tags that have access to the resource
	GetResourceACL(models.Resource) ([]string, error)
	GetResourceACLWithMemberInfo(models.Resource) ([]models.Member, error)
	GetMembersAccess(models.Member) ([]models.MemberAccess, error)
}

type ResourceManager struct {
	mqttServer MQTTServer
	store      RMStore
}

const (
	// StatusGood - the resource is online and up to date
	StatusGood = iota
	// StatusOutOfDate - the resource does not have the most up to date information
	StatusOutOfDate
	// StatusOffline - the resource is not reachable
	StatusOffline
)

func NewResourceManager(ms MQTTServer, store RMStore) *ResourceManager {
	return &ResourceManager{ms, store}
}

// UpdateResourceACL pulls a resource's accesslist from the DB and pushes it to the resource
func (rm *ResourceManager) UpdateResourceACL(r models.Resource) error {
	// get acl for that resource
	accessList, err := rm.store.GetResourceACL(r)

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
	rm.mqttServer.Publish(r.Name+"/update", j)

	return nil
}

// UpdateResources - publish an MQTT message to add a member to the actual device
func (rm *ResourceManager) UpdateResources() {
	resources := rm.store.GetResources()

	for _, r := range resources {
		members, _ := rm.store.GetResourceACLWithMemberInfo(r)
		for _, m := range members {
			b, _ := json.Marshal(&AddMemberRequest{
				ResourceAddress: r.Address,
				Command:         "adduser",
				UserName:        m.Name,
				RFID:            m.RFID,
				AccessType:      1,
				ValidUntil:      -86400,
			})
			rm.mqttServer.Publish(r.Name, string(b))

			time.Sleep(2 * time.Second)
		}
	}
}

// PushOne - update one user on the resources
func (rm *ResourceManager) PushOne(m models.Member) {
	memberAccess, _ := rm.store.GetMembersAccess(m)
	for _, m := range memberAccess {
		b, _ := json.Marshal(&AddMemberRequest{
			ResourceAddress: m.ResourceAddress,
			Command:         "adduser",
			UserName:        m.Name,
			RFID:            m.RFID,
			AccessType:      1,
			ValidUntil:      -86400,
		})
		rm.mqttServer.Publish(m.ResourceName, string(b))
	}
}

func (rm *ResourceManager) DeleteResourceACL() {
	resources := rm.store.GetResources()

	for _, r := range resources {
		b, _ := json.Marshal(&DeleteMemberRequest{
			ResourceAddress: r.Address,
			Command:         "deletusers", // not a type-o this is how the command is defined in the rfid reader
		})
		rm.mqttServer.Publish(r.Name, string(b))
	}
}

// CheckStatus will publish an mqtt command that requests for a specific device to verify that
//   the resource has the correct and up to date access list
//   It will do this by hashing the list retrieved from the DB and comparing it
//   with the hash that the resource reports
func (rm *ResourceManager) CheckStatus(r database.Resource) {
	rm.mqttServer.Publish(r.Name+"/cmd", "aclhash")
}

func hash(accessList []string) string {
	h := sha1.New()
	h.Write([]byte(strings.Join(accessList[:], "\n")))
	bs := h.Sum(nil)

	log.Debug(strings.Join(accessList[:], "\n"))
	log.Debugf("%x\n", bs)
	return fmt.Sprintf("%x", bs)
}
