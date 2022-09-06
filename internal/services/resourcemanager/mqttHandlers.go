package resourcemanager

import (
	"encoding/json"
	"fmt"
	"memberserver/internal/datastore/dbstore"
	"memberserver/internal/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
		rm.logger.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	rm.logger.Infof("name from resource: %s", acl.Name)
	// get resourceByName
	r, err := rm.store.GetResourceByName(acl.Name)
	if err != nil {
		rm.logger.Errorf("error fetching resource: %s", err)
		return
	}
	accessList, err := rm.store.GetResourceACL(r)
	if err != nil {
		rm.logger.Error(err)
		return
	}

	if acl.Hash != rm.hash(accessList) {
		rm.logger.Debugf("[%s] is out of date - attempting to update with new data", r.Name)
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

func (rm *ResourceManager) Receive(client mqtt.Client, msg mqtt.Message) {
	var payload models.LogMessage

	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		rm.logger.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	if payload.EventTime == 0 {
		rm.logger.Println("receive: ", string(msg.Payload()))
		return
	}

	rm.OnAccessEvent(payload)
}

// OnAccessEvent - post the event to slack. This could also get shoved in the DB eventually
func (rm *ResourceManager) OnAccessEvent(payload models.LogMessage) {
	m, err := rm.store.GetMemberByRFID(payload.RFID)
	if err != nil {
		rm.logger.Errorf("swipe on %s of unknown fob: %s", payload.Door, payload.RFID)
		return
	}

	defer func(m models.Member, p models.LogMessage) {
		go rm.notifier.Send(fmt.Sprintf("name: %s, rfid: %s, door: %s, time: %d", m.Name, p.RFID, p.Door, p.EventTime))
		go rm.store.LogAccessEvent(models.LogMessage{
			Type:      p.Type,
			EventTime: p.EventTime,
			IsKnown:   p.IsKnown,
			Username:  m.Name,
			RFID:      p.RFID,
			Door:      p.Door,
		})
	}(m, payload)
}

type HeartBeat struct {
	ResourceName string `json:"door"`
}

// OnHeartBeat handles heartbeats from
func (rm *ResourceManager) OnHeartBeat(client mqtt.Client, msg mqtt.Message) {
	var hb HeartBeat
	err := json.Unmarshal(msg.Payload(), &hb)
	if err != nil {
		rm.logger.Errorf("error unmarshalling mqtt payload: %s", err)
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
