package v1

import (
	"github.com/HackRVA/memberserver/datastore"
	"github.com/HackRVA/memberserver/models"
	mqttserver "github.com/HackRVA/memberserver/pkg/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHandler interface {
	CheckStatus(r models.Resource)
	DeleteResourceACL()
	EnableValidUIDs()
	HealthCheckHandler(client mqtt.Client, msg mqtt.Message)
	MQTT() mqttserver.MQTTServer
	OnAccessEventHandler(payload models.LogMessage)
	OnHeartBeatHandler(client mqtt.Client, msg mqtt.Message)
	OnRemoveInvalidRequestHandler(client mqtt.Client, msg mqtt.Message)
	Open(resource models.Resource)
	PushOne(m models.Member)
	ReceiveHandler(client mqtt.Client, msg mqtt.Message)
	RemoveMember(memberAccess models.MemberAccess)
	RemoveOne(member models.Member)
	RemovedInvalidUIDs()
	UpdateResourceACL(r models.Resource) error
}

const (
	commandDeleteUID = "deletuid"
	commandAddUser   = "adduser"
	commandOpenDoor  = "opendoor"
	commandListUser  = "listusr"
)

type notifier interface {
	Send(msg string)
}

type mqttServer interface {
	Publish(address string, topic string, payload interface{})
	Subscribe(address string, topic string, handler mqtt.MessageHandler)
}

type mqttHandler struct {
	datastore.DataStore
	mqtt     mqttserver.MQTTServer
	notifier notifier
}

func New(ms mqttServer, store datastore.DataStore, notifier notifier) *mqttHandler {
	return &mqttHandler{mqtt: ms, DataStore: store, notifier: notifier}
}

func (v1 mqttHandler) MQTT() mqttserver.MQTTServer {
	return v1.mqtt
}
