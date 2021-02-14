package resourcemanager

import (
	"encoding/json"
	"fmt"
	"memberserver/database"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// healthCheck
//  this is the mqtt messageHandler that runs when a resource checks in
//  we expect the payload to be json that marshals to `ACLResponse` which includes the name
//  and a hash of it's ACL
//  if the ACL hash doesn't match what we have in the database, we will trigger an update to push
//  to the resource
var healthCheck mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
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

// Subscribe - subscribe to an MQTT topic and pass in a messageHandler
func Subscribe(topic string, handler mqtt.MessageHandler) {
	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER_ADDRESS"))
	opts.SetClientID("member-server")

	opts.OnConnect = func(c mqtt.Client) {
		if token := c.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
			log.Error(token.Error())
		}
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
		return
	}
	log.Debug("Connected to server\n")
}

// Publish - publish to an MQTT topic
func Publish(topic string, payload interface{}) {
	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER_ADDRESS"))
	opts.SetClientID("member-server")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
		return
	}

	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	client.Disconnect(250)
}
