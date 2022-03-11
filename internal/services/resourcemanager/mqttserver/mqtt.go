package mqttserver

import (
	"math/rand"
	"memberserver/internal/services/config"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MQTTServer interface {
	Publish(topic string, payload interface{})
	Subscribe(topic string, handler mqtt.MessageHandler)
}

type mqttServer struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewMQTTServer() *mqttServer {
	return &mqttServer{}
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
func (m *mqttServer) Subscribe(topic string, handler mqtt.MessageHandler) {
	conf, _ := config.Load()
	opts := mqtt.NewClientOptions().AddBroker(conf.MQTTBrokerAddress)
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
func (m *mqttServer) Publish(topic string, payload interface{}) {
	conf, _ := config.Load()
	opts := mqtt.NewClientOptions().AddBroker(conf.MQTTBrokerAddress)
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
