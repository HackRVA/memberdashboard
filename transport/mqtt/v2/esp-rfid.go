package v1

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HackRVA/memberserver/datastore/dbstore"
	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/services/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// HealthCheck -- this is the mqtt messageHandler that runs when a resource checks in
//
//	we expect the payload to be json that marshals to `ACLResponse` which includes the name
//	and a hash of it's ACL
//	if the ACL hash doesn't match what we have in the database, we will trigger an update to push
//	to the resource
func (v1 mqttHandler) HealthCheckHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())

	var acl models.ACLResponse

	err := json.Unmarshal(msg.Payload(), &acl)
	if err != nil {
		logger.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	logger.Infof("name from resource: %s", acl.Name)
	// get resourceByName
	r, err := v1.GetResourceByName(acl.Name)
	if err != nil {
		logger.Errorf("error fetching resource: %s", err)
		return
	}
	accessList, err := v1.GetResourceACL(r)
	if err != nil {
		logger.Error(err)
		return
	}

	if acl.Hash != hash(accessList) {
		logger.Debugf("[%s] is out of date - attempting to update with new data", r.Name)
		// status = StatusOutOfDate
		// err = UpdateResourceACL(r)
		// if err != nil {
		// 	log.Errorf("error updating resource with acl: %s", err)
		// }
	}

	// TODO: check that the resource responds with a hash of the list
	// status = StatusGood
}

func hash(accessList []string) string {
	h := sha1.New()
	h.Write([]byte(strings.Join(accessList[:], "\n")))
	bs := h.Sum(nil)

	logger.Debug(strings.Join(accessList[:], "\n"))
	logger.Debugf("%x\n", bs)
	return fmt.Sprintf("%x", bs)
}

type HeartBeat struct {
	Type         string `json:"type"`
	Time         int64  `json:"time"`
	Uptime       int64  `json:"uptime"`
	IP           string `json:"ip"`
	ResourceName string `json:"hostname"`
}

// OnHeartBeat handles heartbeats from
func (v1 mqttHandler) OnHeartBeatHandler(client mqtt.Client, msg mqtt.Message) {
	var hb HeartBeat
	err := json.Unmarshal(msg.Payload(), &hb)
	if err != nil {
		logger.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	dbstore.ResourceHeartbeat(models.Resource{
		Name: hb.ResourceName,
	})
}
