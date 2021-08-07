package resourcemanager

import (
	"encoding/base64"
	"encoding/json"
	"memberserver/api/models"
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var pub string

type stubMQTTServer struct{}

func (mqtt *stubMQTTServer) Publish(topic string, payload interface{}) {
	json, _ := json.Marshal(payload)
	pub = topic + string(json)
}

func (mqtt *stubMQTTServer) Subscribe(topic string, handler mqtt.MessageHandler) {

}

type stubRMStore struct{}

func (s *stubRMStore) GetResources() []models.Resource {
	return []models.Resource{}
}

func (s *stubRMStore) GetResourceACL(models.Resource) ([]string, error) {
	return []string{}, nil
}

func (s *stubRMStore) GetResourceACLWithMemberInfo(models.Resource) ([]models.Member, error) {
	return []models.Member{}, nil
}

func (s *stubRMStore) GetMembersAccess(models.Member) ([]models.MemberAccess, error) {
	return []models.MemberAccess{}, nil
}

func TestUpdateResourceACL(t *testing.T) {
	resourceManager := NewResourceManager(&stubMQTTServer{}, &stubRMStore{})
	resource := models.Resource{
		ID:   "0",
		Name: "should just straight up send it",
	}

	// ACLUpdateRequest{}

	acl := `{"acl":[]}`
	// byte, _ := json.Marshal(resource.Name + "/update")
	want := "should just straight up send it/update\"" + base64.RawStdEncoding.EncodeToString([]byte(acl)) + "==\""

	resourceManager.UpdateResourceACL(resource)
	if pub != want {
		t.Errorf("did not succeed. got %s want: %s", pub, want)
	}
}
