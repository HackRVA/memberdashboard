package resourcemanager

import (
	"encoding/json"
	"fmt"
	"memberserver/internal/datastore/dbstore"
	"memberserver/internal/models"
	"memberserver/pkg/slack"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// HealthCheck -- this is the mqtt messageHandler that runs when a resource checks in
//  we expect the payload to be json that marshals to `ACLResponse` which includes the name
//  and a hash of it's ACL
//  if the ACL hash doesn't match what we have in the database, we will trigger an update to push
//  to the resource
func (rm *ResourceManager) HealthCheck(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())

	var acl models.ACLResponse

	err := json.Unmarshal(msg.Payload(), &acl)
	if err != nil {
		log.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	log.Infof("name from resource: %s", acl.Name)
	// get resourceByName
	r, err := rm.store.GetResourceByName(acl.Name)
	if err != nil {
		log.Errorf("error fetching resource: %s", err)
		return
	}
	accessList, err := rm.store.GetResourceACL(r)
	if err != nil {
		log.Error(err)
		return
	}

	if acl.Hash != hash(accessList) {
		log.Debugf("[%s] is out of date - attempting to update with new data", r.Name)
		// status = StatusOutOfDate
		// err = UpdateResourceACL(r)
		// if err != nil {
		// 	log.Errorf("error updating resource with acl: %s", err)
		// }
	}

	// TODO: check that the resource responds with a hash of the list
	// status = StatusGood
}

// {"cmd":"log","type":"access","time":1631240207,"isKnown":"true","access":"Always","username":"Stanley Hash","uid":"f3ec6234","door":"frontdoor"}
type EventLogPayload struct {
	Time     int    `json:"time"`
	Username string `json:"username"`
	RFID     string `json:"uid"`
	Door     string `json:"door"`
}

// OnAccessEvent - post the event to slack. This could also get shoved in the DB eventually
func (rm *ResourceManager) OnAccessEvent(client mqtt.Client, msg mqtt.Message) {
	var payload models.LogMessage

	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	log.Println(string(msg.Payload()))
	go slack.Send(fmt.Sprintf("name: %s, rfid: %s, door: %s, time: %d", payload.Username, payload.RFID, payload.Door, payload.EventTime))
	go rm.store.LogAccessEvent(models.LogMessage{
		Type:      payload.Type,
		EventTime: payload.EventTime,
		IsKnown:   payload.IsKnown,
		Username:  payload.Username,
		RFID:      payload.RFID,
		Door:      payload.Door,
	})
}

type HeartBeat struct {
	ResourceName string `json:"door"`
}

// OnHeartBeat handles heartbeats from
func (rm *ResourceManager) OnHeartBeat(client mqtt.Client, msg mqtt.Message) {
	var hb HeartBeat
	err := json.Unmarshal(msg.Payload(), &hb)
	if err != nil {
		log.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	dbstore.ResourceHeartbeat(models.Resource{
		Name: hb.ResourceName,
	})
}

// go through and remove members rfid fobs that are listed as invalid
func (rm *ResourceManager) OnRemoveInvalidRequest(client mqtt.Client, msg mqtt.Message) {
	rm.RemovedInvalidUIDs()
}
