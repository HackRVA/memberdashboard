package resourcemanager

import (
	"encoding/json"
	"fmt"
	"memberserver/database"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// MQTTClient - how we connect to the MQTT Broker
type MQTTClient struct {
	opts   *mqtt.ClientOptions
	client mqtt.Client
}

var check mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
	// status := StatusOffline

	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}

	var acl ACLResponse

	_ = json.Unmarshal(msg.Payload(), &acl)
	// get resourceByName
	r, err := db.GetResourceByName(acl.Name)
	accessList, err := db.GetResourceACL(r)

	log.Debugf("body= %s json=%s accessListHash=%s", string(msg.Payload()), acl.Hash, hash(accessList))
	if acl.Hash != hash(accessList) {
		log.Debugf("attempting to update resource [%s] with new data", r.Name)
		// status = StatusOutOfDate
		// err = UpdateResourceACL(r)
		if err != nil {
			log.Errorf("error updating resource with acl: %s", err)
		}
	}

	// TODO: check that the resource responds with a hash of the list
	// status = StatusGood

	db.Release()
}

// MQTTSetup setup mqtt client
func MQTTSetup() MQTTClient {
	var m MQTTClient
	m.opts = mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER_ADDRESS"))
	m.opts.SetClientID("member-server")
	m.opts.SetDefaultPublishHandler(check)

	return m
}

// Subscribe - subscribe to an MQTT topic and pass in a messageHandler
func (m MQTTClient) Subscribe(topic string, handler mqtt.MessageHandler) {
	m.opts.OnConnect = func(c mqtt.Client) {
		if token := c.Subscribe(topic, 0, handler); token.Wait() && token.Error() != nil {
			log.Error(token.Error())
		}
	}
	m.client = mqtt.NewClient(m.opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	} else {
		log.Error("Connected to server\n")
	}
}

// Publish - publish to an MQTT topic
func (m MQTTClient) Publish(topic string, payload interface{}) {
	token := m.client.Publish(topic, 0, false, payload)
	token.Wait()
}
