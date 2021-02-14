package resourcemanager

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"memberserver/database"
)

// Resource manager keeps the resources up to date by
//  pushing new updates and checking in on their health

// statusCheckInterval - check the resources every hour
const statusCheckInterval = 1

// ResourceManager contains functions that
type ResourceManager struct {
	db *database.Database
	mc MQTTClient
}

// ACLUpdateRequest is the json object we send to a resource when pushing an update
type ACLUpdateRequest struct {
	ACL []string `json:"acl"`
}

const (
	// StatusGood - the resource is online and up to date
	StatusGood = iota
	// StatusOutOfDate - the resource does not have the most up to date information
	StatusOutOfDate
	// StatusOffline - the resource is not reachable
	StatusOffline
)

// Setup initializes the resource manager
func Setup() (*ResourceManager, error) {
	var err error
	rm := &ResourceManager{}
	rm.db, err = database.Setup()

	if err != nil {
		log.Errorf("error setting up db: %s", err)
		return rm, err
	}
	rm.mc = MQTTSetup()

	rm.mc.Subscribe("/result", check)

	// quietly check the resource status on an interval
	ticker := time.NewTicker(statusCheckInterval * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				resources := rm.db.GetResources()

				for _, r := range resources {
					rm.CheckStatus(r)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return rm, err
}

// UpdateResourceACL pulls a resource's accesslist from the DB and pushes it to the resource
func (rm ResourceManager) UpdateResourceACL(r database.Resource) error {
	// get acl for that resource
	accessList, err := rm.db.GetResourceACL(r)

	if err != nil {
		return err
	}

	updateRequest := &ACLUpdateRequest{}
	updateRequest.ACL = accessList

	j, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}
	log.Debugf("access list: %s", j)

	// publish the update to mqtt broker
	rm.mc.Publish(r.Name+"/update", j)

	return nil
}

// ACLResponse Response from a resource that is a hash of the ACL that the
//   resource has stored
type ACLResponse struct {
	Hash string `json:"acl"`
	// Name of the resource - this should match what we have in the database
	//  so we know which acl to compare it with
	Name string `json:"name"`
}

// CheckStatus will publish an mqtt command that requests for a specific device to verify that
//   the resource has the correct and up to date access list
//   It will do this by hashing the list retrieved from the DB and comparing it
//   with the hash that the resource reports
func (rm ResourceManager) CheckStatus(r database.Resource) (uint8, error) {
	var status uint8
	accessList, err := rm.db.GetResourceACL(r)
	status = StatusOffline

	if err != nil {
		return status, err
	}

	// push the update to the resource
	resp, err := http.Get(r.Address)
	if err != nil {
		log.Errorf("Unable to reach the resource.")
		return status, err
	}
	var acl ACLResponse
	body, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(body, &acl)

	log.Debugf("body= %s json=%s accessListHash=%s", string(body), acl.Hash, hash(accessList))
	if acl.Hash != hash(accessList) {
		log.Debugf("attempting to update resource [%s] with new data", r.Name)
		err = rm.UpdateResourceACL(r)
		status = StatusOutOfDate
		if err != nil {
			log.Errorf("error updating resource with acl: %s", err)
			return status, err
		}
		return status, err
	}

	// TODO: check that the resource responds with a hash of the list
	status = StatusGood

	return status, nil
}

func hash(accessList []string) string {
	h := sha1.New()
	h.Write([]byte(strings.Join(accessList[:], "\n")))
	bs := h.Sum(nil)

	log.Debug(strings.Join(accessList[:], "\n"))
	log.Debugf("%x\n", bs)
	return fmt.Sprintf("%x", bs)
}
