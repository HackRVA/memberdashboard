package resource

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"memberserver/internal/datastore"
	"memberserver/internal/models"
	"memberserver/internal/services/logger"
	"memberserver/pkg/mqtt"
	"time"

	"strings"

	"github.com/sirupsen/logrus"
)

const (
	commandDeleteUID = "deletuid"
	commandAddUser   = "adduser"
	commandOpenDoor  = "opendoor"
)

// Resource manager keeps the resources up to date by
//  pushing new updates and checking in on their health

type Resource struct {
	MQTTServer mqtt.MQTTServer
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

func NewResource(ms mqtt.MQTTServer, store datastore.DataStore) Resource {
	return Resource{ms, store}
}

// UpdateResourceACL pulls a resource's accesslist from the DB and pushes it to the resource
func (rm Resource) UpdateResourceACL(r models.Resource) error {
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
	logger.Infof("access list: %s", j)

	// publish the update to mqtt broker
	rm.MQTTServer.Publish(r.Name+"/update", j)

	return nil
}

func (rm Resource) GetResources() []models.Resource {
	return rm.store.GetResources()
}

func (rm Resource) GetOne(id string) (models.Resource, error) {
	return rm.store.GetResourceByID(id)
}

// Update - update properties of a resource
func (r Resource) Update(res models.Resource) (*models.Resource, error) {
	return r.store.UpdateResource(res)
}

func (r Resource) Delete(res models.Resource) error {
	return r.store.DeleteResource(res.ID)
}

// UpdateResources - publish an MQTT message to add a member to the actual device
func (rm Resource) UpdateResources() {
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

func (rm Resource) EnableValidUIDs() {
	activeMembers, err := rm.store.GetActiveMembersByResource()
	if err != nil {
		logger.Errorf("error getting active members from db %s", err.Error())
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

func (rm Resource) RemovedInvalidUIDs() {
	inactiveMembers, err := rm.store.GetInactiveMembersByResource()
	if err != nil {
		logger.Errorf("error getting inactive members from db %s", err.Error())
		return
	}

	logger.Debug("looking for members to remove")

	for _, m := range inactiveMembers {
		/* We will just try to remove all invalid members even if they are already removed */
		rm.RemoveMember(m)
		time.Sleep(2 * time.Second)
	}
}

func (rm Resource) RemoveMember(memberAccess models.MemberAccess) {
	b, _ := json.Marshal(&models.MemberRequest{
		ResourceAddress: memberAccess.ResourceAddress,
		Command:         commandDeleteUID,
		RFID:            memberAccess.RFID,
	})
	rm.MQTTServer.Publish(memberAccess.ResourceName, string(b))
	logger.Debugf("attempting to remove member %s from rfid device %s : %s", memberAccess.Email, memberAccess.ResourceName, memberAccess.ResourceAddress)
}

func (rm Resource) Open(resource models.Resource) {
	b, _ := json.Marshal(models.MQTTRequest{
		Door:    resource.Name,
		Command: commandOpenDoor,
		Address: resource.Address,
	})

	rm.MQTTServer.Publish(resource.Name, string(b))
}

func (rm Resource) RemoveOneByResourceID(member models.Member, resourceID string) error {
	r, err := rm.GetOne(resourceID)
	if err != nil {
		return err
	}

	rm.store.RemoveUserFromResource(member.Email, r.Name)
	rm.RemoveMember(models.MemberAccess{
		Email:           member.Email,
		ResourceAddress: r.Address,
		ResourceName:    r.Name,
		Name:            member.Name,
		RFID:            member.RFID,
	})
	return nil
}

// RemoveOne - remove a member from all resources
func (rm Resource) RemoveOne(member models.Member) {
	member, err := rm.store.GetMemberByEmail(member.Email)
	if err != nil {
		logrus.Error(err)
		return
	}

	memberAccess, _ := rm.store.GetMembersAccess(member)

	for _, m := range memberAccess {
		rm.store.RemoveUserFromResource(member.Email, m.ResourceName)
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
func (rm Resource) PushOne(m models.Member) {
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

func (rm Resource) DeleteResourceACL() {
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
func (rm Resource) CheckStatus(r models.Resource) {
	rm.MQTTServer.Publish(r.Name+"/cmd", "aclhash")
}

func hash(accessList []string) string {
	h := sha1.New()
	h.Write([]byte(strings.Join(accessList[:], "\n")))
	bs := h.Sum(nil)

	logger.Debug(strings.Join(accessList[:], "\n"))
	logger.Debugf("%x\n", bs)
	return fmt.Sprintf("%x", bs)
}
