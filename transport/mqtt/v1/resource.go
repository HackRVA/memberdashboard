package v1

import (
	"encoding/json"
	"time"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/services/logger"
)

func (v1 mqttHandler) Open(resource models.Resource) {
	b, _ := json.Marshal(models.MQTTRequest{
		Door:    resource.Name,
		Command: commandOpenDoor,
		Address: resource.Address,
	})

	v1.mqtt.Publish(config.Get().MQTTBrokerAddress, resource.Name, string(b))
}

// UpdateResourceACL pulls a resource's accesslist from the DB and pushes it to the resource
func (v1 mqttHandler) UpdateResourceACL(r models.Resource) error {
	// get acl for that resource
	accessList, err := v1.GetResourceACL(r)
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
	v1.mqtt.Publish(config.Get().MQTTBrokerAddress, r.Name+"/update", j)

	return nil
}

func (v1 mqttHandler) DeleteResourceACL() {
	resources := v1.GetResources()

	for _, r := range resources {
		b, _ := json.Marshal(&models.DeleteMemberRequest{
			ResourceAddress: r.Address,
			Command:         "deletusers", // not a type-o this is how the command is defined in the rfid reader
		})
		v1.mqtt.Publish(config.Get().MQTTBrokerAddress, r.Name, string(b))
	}
}

// CheckStatus will publish an mqtt command that requests for a specific device to verify that
//
//	the resource has the correct and up to date access list
//	It will do this by hashing the list retrieved from the DB and comparing it
//	with the hash that the resource reports
func (v1 mqttHandler) CheckStatus(r models.Resource) {
	v1.mqtt.Publish(config.Get().MQTTBrokerAddress, r.Name+"/cmd", "aclhash")
}

// UpdateResources - publish an MQTT message to add a member to the actual device
func (v1 mqttHandler) UpdateResources() {
	resources := v1.GetResources()

	for _, r := range resources {
		members, _ := v1.GetResourceACLWithMemberInfo(r)
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
			v1.mqtt.Publish(config.Get().MQTTBrokerAddress, r.Name, string(b))

			time.Sleep(2 * time.Second)
		}
	}
}
