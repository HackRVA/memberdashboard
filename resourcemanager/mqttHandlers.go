package resourcemanager

import (
	"encoding/json"
	"fmt"
	"memberserver/api/models"
	"memberserver/datastore/dbstore"
	"memberserver/slack"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// HealthCheck -- this is the mqtt messageHandler that runs when a resource checks in
//  we expect the payload to be json that marshals to `ACLResponse` which includes the name
//  and a hash of it's ACL
//  if the ACL hash doesn't match what we have in the database, we will trigger an update to push
//  to the resource
var HealthCheck mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())

	var acl models.ACLResponse

	err := json.Unmarshal(msg.Payload(), &acl)
	if err != nil {
		log.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	log.Debugf("name from resource: %s", acl.Name)
	// get resourceByName
	r, err := db.GetResourceByName(acl.Name)
	if err != nil {
		log.Errorf("error fetching resource: %s", err)
		return
	}
	accessList, err := db.GetResourceACL(r)

	// log.Debugf("body= %s json=%s accessListHash=%s name=%s", string(msg.Payload()), acl.Hash, hash(accessList), acl.Name)

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
var OnAccessEvent mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var payload EventLogPayload

	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		log.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	member, err := db.GetMemberByRFID(payload.RFID)
	if err != nil {
		log.Errorf("error looking up members name: %s %s", err, string(msg.Payload()))
	}

	slack.PostWebHook(fmt.Sprintf("name: %s, rfid: %s, door: %s, time: %s", member.Name, payload.RFID, payload.Door, payload.Time))

	err = db.AddLogMsg(msg.Payload())
	if err != nil {
		log.Errorf("error saving access event: %s %s", err, string(msg.Payload()))
	}
}

type HeartBeat struct {
	ResourceName string `json:"door"`
}

// OnHeartBeat handles heartbeats from
var OnHeartBeat mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
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
