package mqtt

import (
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MQTTServer interface {
	Publish(address string, topic string, payload interface{})
	Subscribe(address string, topic string, handler mqtt.MessageHandler)
}

type server struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New() *server {
	return &server{}
}

func randStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Subscribe - subscribe to an MQTT topic and pass in a messageHandler
func (m *server) Subscribe(address string, topic string, handler mqtt.MessageHandler) {
	if address == "" {
		log.Error("mqtt address isn't configured")
		return
	}
	opts := mqtt.NewClientOptions().AddBroker(address)
	opts.SetClientID("member-server-subscriber-" + randStringRunes(12))

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
func (m *server) Publish(address string, topic string, payload interface{}) {
	if address == "" {
		log.Error("mqtt address isn't configured")
		return
	}
	opts := mqtt.NewClientOptions().AddBroker(address)
	opts.SetClientID("member-server-publisher-" + randStringRunes(12))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
		return
	}

	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	client.Disconnect(250)
}
