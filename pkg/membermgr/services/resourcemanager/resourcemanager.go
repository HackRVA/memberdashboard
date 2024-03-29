package resourcemanager

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/pkg/membermgr/datastore"
	"github.com/HackRVA/memberserver/pkg/membermgr/models"
	"github.com/HackRVA/memberserver/pkg/mqtt"

	"time"

	"strings"
)

const (
	commandDeleteUID = "deletuid"
	commandAddUser   = "adduser"
	commandOpenDoor  = "opendoor"
	commandListUser  = "listusr"
)

// Resource manager keeps the resources up to date by
//  pushing new updates and checking in on their health

type ResourceManager struct {
	mqtt mqtt.MQTTServer
	datastore.DataStore
	notifier notifier
	logger   logger
}

const (
	// StatusGood - the resource is online and up to date
	StatusGood = iota
	// StatusOutOfDate - the resource does not have the most up to date information
	StatusOutOfDate
	// StatusOffline - the resource is not reachable
	StatusOffline
)

func New(ms mqttServer, store datastore.DataStore, notifier notifier, logger logger) *ResourceManager {
	return &ResourceManager{ms, store, notifier, logger}
}

func (rm ResourceManager) MQTT() mqtt.MQTTServer {
	return rm.mqtt
}

// UpdateResourceACL pulls a resource's accesslist from the DB and pushes it to the resource
func (rm ResourceManager) UpdateResourceACL(r models.Resource) error {
	// get acl for that resource
	accessList, err := rm.GetResourceACL(r)

	if err != nil {
		return err
	}

	updateRequest := &models.ACLUpdateRequest{}
	updateRequest.ACL = accessList

	j, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}
	rm.logger.Infof("access list: %s", j)

	// publish the update to mqtt broker
	rm.mqtt.Publish(config.Get().MQTTBrokerAddress, r.Name+"/update", j)

	return nil
}

// UpdateResources - publish an MQTT message to add a member to the actual device
func (rm ResourceManager) UpdateResources() {
	resources := rm.GetResources()

	for _, r := range resources {
		members, _ := rm.GetResourceACLWithMemberInfo(r)
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
			rm.mqtt.Publish(config.Get().MQTTBrokerAddress, r.Name, string(b))

			time.Sleep(2 * time.Second)
		}
	}
}

func (rm ResourceManager) EnableValidUIDs() {
	activeMembers, err := rm.GetActiveMembersByResource()
	if err != nil {
		rm.logger.Errorf("error getting active members from db %s", err.Error())
		return
	}

	for _, m := range activeMembers {
		rm.PushOne(models.Member{
			Name: m.Name,
			RFID: m.RFID,
		})
		time.Sleep(2 * time.Second)
	}
}

func (rm ResourceManager) RemovedInvalidUIDs() {
	inactiveMembers, err := rm.GetInactiveMembersByResource()
	if err != nil {
		rm.logger.Errorf("error getting inactive members from db %s", err.Error())
		return
	}

	rm.logger.Debug("looking for members to remove")

	for _, m := range inactiveMembers {
		/* We will just try to remove all invalid members even if they are already removed */
		rm.RemoveMember(m)
		time.Sleep(2 * time.Second)
	}
}

func (rm ResourceManager) RemoveMember(memberAccess models.MemberAccess) {
	b, _ := json.Marshal(&models.MemberRequest{
		ResourceAddress: memberAccess.ResourceAddress,
		Command:         commandDeleteUID,
		RFID:            memberAccess.RFID,
	})

	rm.mqtt.Publish(config.Get().MQTTBrokerAddress, memberAccess.ResourceName, string(b))
	rm.logger.Debugf("attempting to remove member %s from rfid device %s : %s", memberAccess.Email, memberAccess.ResourceName, memberAccess.ResourceAddress)
}

func (rm ResourceManager) Open(resource models.Resource) {
	b, _ := json.Marshal(models.MQTTRequest{
		Door:    resource.Name,
		Command: commandOpenDoor,
		Address: resource.Address,
	})

	rm.mqtt.Publish(config.Get().MQTTBrokerAddress, resource.Name, string(b))
}

// RemoveOne - remove a member from all resources
func (rm ResourceManager) RemoveOne(member models.Member) {
	member, err := rm.GetMemberByEmail(member.Email)
	if err != nil {
		rm.logger.Error(err)
		return
	}

	memberAccess, _ := rm.GetMembersAccess(member)

	for _, m := range memberAccess {
		rm.RemoveMember(models.MemberAccess{
			Email:           member.Email,
			ResourceAddress: m.ResourceAddress,
			ResourceName:    m.ResourceName,
			Name:            member.Name,
			RFID:            member.RFID,
		})
	}
}

// PushOne - update one user on the resources
func (rm ResourceManager) PushOne(m models.Member) {
	memberAccess, _ := rm.GetMembersAccess(m)
	for _, m := range memberAccess {
		b, _ := json.Marshal(&models.MemberRequest{
			ResourceAddress: m.ResourceAddress,
			Command:         commandAddUser,
			UserName:        m.Name,
			RFID:            m.RFID,
			AccessType:      1,
			ValidUntil:      -86400,
		})
		rm.mqtt.Publish(config.Get().MQTTBrokerAddress, m.ResourceName, string(b))
	}
}

func (rm ResourceManager) DeleteResourceACL() {
	resources := rm.GetResources()

	for _, r := range resources {
		b, _ := json.Marshal(&models.DeleteMemberRequest{
			ResourceAddress: r.Address,
			Command:         "deletusers", // not a type-o this is how the command is defined in the rfid reader
		})
		rm.mqtt.Publish(config.Get().MQTTBrokerAddress, r.Name, string(b))
	}
}

// CheckStatus will publish an mqtt command that requests for a specific device to verify that
//
//	the resource has the correct and up to date access list
//	It will do this by hashing the list retrieved from the DB and comparing it
//	with the hash that the resource reports
func (rm ResourceManager) CheckStatus(r models.Resource) {
	rm.mqtt.Publish(config.Get().MQTTBrokerAddress, r.Name+"/cmd", "aclhash")
}

func (rm ResourceManager) hash(accessList []string) string {
	h := sha1.New()
	h.Write([]byte(strings.Join(accessList[:], "\n")))
	bs := h.Sum(nil)

	rm.logger.Debug(strings.Join(accessList[:], "\n"))
	rm.logger.Debugf("%x\n", bs)
	return fmt.Sprintf("%x", bs)
}
