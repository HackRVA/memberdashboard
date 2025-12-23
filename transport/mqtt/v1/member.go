package v1

import (
	"encoding/json"
	"time"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/services/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (v1 mqttHandler) EnableValidUIDs() {
	activeMembers, err := v1.GetActiveMembersByResource()
	if err != nil {
		logger.Errorf("error getting active members from db %s", err.Error())
		return
	}

	for _, m := range activeMembers {
		v1.PushOne(models.Member{
			Name: m.Name,
			RFID: m.RFID,
		})
		time.Sleep(2 * time.Second)
	}
}

func (v1 mqttHandler) RemovedInvalidUIDs() {
	inactiveMembers, err := v1.GetInactiveMembersByResource()
	if err != nil {
		logger.Errorf("error getting inactive members from db %s", err.Error())
		return
	}

	logger.Debug("looking for members to remove")

	for _, m := range inactiveMembers {
		/* We will just try to remove all invalid members even if they are already removed */
		v1.RemoveMember(m)
		time.Sleep(2 * time.Second)
	}
}

func (v1 mqttHandler) RemoveMember(memberAccess models.MemberAccess) {
	b, _ := json.Marshal(&models.MemberRequest{
		ResourceAddress: memberAccess.ResourceAddress,
		Command:         commandDeleteUID,
		RFID:            memberAccess.RFID,
	})

	v1.mqtt.Publish(config.Get().MQTTBrokerAddress, memberAccess.ResourceName+"/cmd", string(b))
	logger.Debugf("attempting to remove member %s from rfid device %s : %s", memberAccess.Email, memberAccess.ResourceName, memberAccess.ResourceAddress)
}

// go through and remove members rfid fobs that are listed as invalid
func (v1 *mqttHandler) OnRemoveInvalidRequestHandler(client mqtt.Client, msg mqtt.Message) {
	v1.RemovedInvalidUIDs()
}

// PushOne - update one user on the resources
func (v1 mqttHandler) PushOne(m models.Member) {
	memberAccess, _ := v1.GetMembersAccess(m)
	for _, m := range memberAccess {
		b, _ := json.Marshal(&models.MemberRequest{
			ResourceAddress: m.ResourceAddress,
			Command:         commandAddUser,
			UserName:        m.Name,
			RFID:            m.RFID,
			AccessType:      1,
			ValidUntil:      -86400,
		})
		v1.mqtt.Publish(config.Get().MQTTBrokerAddress, m.ResourceName+"/cmd", string(b))
	}
}

// RemoveOne - remove a member from all resources
func (v1 mqttHandler) RemoveOne(member models.Member) {
	member, err := v1.GetMemberByEmail(member.Email)
	if err != nil {
		logger.Error(err)
		return
	}

	memberAccess, _ := v1.GetMembersAccess(member)

	for _, m := range memberAccess {
		v1.RemoveMember(models.MemberAccess{
			Email:           member.Email,
			ResourceAddress: m.ResourceAddress,
			ResourceName:    m.ResourceName,
			Name:            member.Name,
			RFID:            member.RFID,
		})
	}
}
