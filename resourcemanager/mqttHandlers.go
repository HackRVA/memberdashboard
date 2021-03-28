package resourcemanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"memberserver/config"
	"memberserver/database"
	"net/http"

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
	// status := StatusOffline

	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	var acl ACLResponse

	err = json.Unmarshal(msg.Payload(), &acl)
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
		err = UpdateResourceACL(r)
		if err != nil {
			log.Errorf("error updating resource with acl: %s", err)
		}
	}

	// TODO: check that the resource responds with a hash of the list
	// status = StatusGood

	db.Release()
}

// OnAccessEvent - post the event to slack. This could also get shoved in the DB eventually
var OnAccessEvent mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	conf, _ := config.Load()

	newMsg := fmt.Sprint("{\"text\":'```", string(msg.Payload()), "```'}")
	jsonStr := []byte(newMsg)
	log.Debugf("attempting to post to slack %s", newMsg)

	c := &http.Client{}
	req, err := http.NewRequest("POST", conf.SlackAccessEvents, bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Errorf("some error sending to slack hook: %s", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.Do(req)

	if err != nil {
		log.Errorf("some error sending to slack hook: %s", err)
		return
	}

	defer res.Body.Close()
}
