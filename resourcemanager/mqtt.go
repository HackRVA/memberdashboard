package resourcemanager

import (
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

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
