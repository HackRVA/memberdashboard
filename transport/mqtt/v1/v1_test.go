package v1_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/HackRVA/memberserver/datastore/in_memory"
	v1 "github.com/HackRVA/memberserver/transport/mqtt/v1"

	"github.com/HackRVA/memberserver/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var pub []string

type slackNotifier struct{}

func (s slackNotifier) Send(msg string) {}

type stubMQTTServer struct{}

func (mqtt *stubMQTTServer) Publish(address string, topic string, payload interface{}) {
	println(topic)
	json, _ := json.Marshal(payload)
	pub = append(pub, topic+string(json))
}

func (mqtt *stubMQTTServer) Subscribe(address string, topic string, handler mqtt.MessageHandler) {
}

// TestUpdateResourceACL we just want to test that the mqtt message looks reasonable
func TestUpdateResourceACL(t *testing.T) {
	resourceManager := v1.New(&stubMQTTServer{}, &in_memory.In_memory{}, slackNotifier{})
	resource := models.Resource{
		ID:   "0",
		Name: "should just straight up send it",
	}

	acl := `{"acl":[]}`
	want := "should just straight up send it/update\"" + base64.RawStdEncoding.EncodeToString([]byte(acl)) + "==\""

	if err := resourceManager.UpdateResourceACL(resource); err != nil {
		t.Errorf("error updating resource acl %s", err)
	}
	if pub[0] != want {
		t.Errorf("did not succeed. got: %s want: %s", pub[0], want)
	}
	pub = []string{}
}

// TestUpdateResources we just want to test that the mqtt message looks reasonable
func TestUpdateResources(t *testing.T) {
	resourceManager := v1.New(&stubMQTTServer{}, &in_memory.In_memory{}, slackNotifier{})
	resources := []models.Resource{
		{
			ID:   "0",
			Name: "should just straight up send it",
		},
		{
			ID:   "1",
			Name: "1should just straight up send it",
		},
		{
			ID:   "2",
			Name: "2should just straight up send it",
		},
	}

	// add some stuff to the store
	for _, v := range resources {
		if _, err := resourceManager.RegisterResource(v.Name, v.Address, v.IsDefault); err != nil {
			t.Errorf("error registering resource %s", err)
		}
	}

	want := `should just straight up send it"{\"doorip\":\"\",\"cmd\":\"adduser\",\"user\":\"test\",\"uid\":\"\",\"acctype\":1,\"validuntil\":-86400}"`
	resourceManager.UpdateResources()
	if len(pub) != 3 {
		t.Errorf("it didn't send all of the updates, received: %d", len(pub))
	}

	for i, v := range resources {
		if v != resources[i] {
			t.Errorf("does not look like a proper resource. got: %s want: %s", pub, want)
		}
	}

	pub = []string{}
}
