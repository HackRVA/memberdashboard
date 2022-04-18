package jobs

import (
	"io/ioutil"
	"memberserver/internal/datastore"
	"memberserver/internal/integrations"
	"memberserver/internal/models"
	"memberserver/internal/services/config"
	"memberserver/internal/services/logger"
	"memberserver/internal/services/mail"
	"memberserver/internal/services/member"
	"memberserver/internal/services/resourcemanager"
	"memberserver/internal/services/resourcemanager/mqttserver"
	"memberserver/pkg/paypal"

	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type JobController struct {
	config          config.Config
	DataStore       datastore.DataStore
	mailAPI         mail.MailApi
	resourceManager resourcemanager.ResourceManager
	paymentProvider integrations.PaymentProvider
	member          member.MemberService
}

func New(db datastore.DataStore) JobController {
	config, _ := config.Load()
	mailAPI, _ := mail.Setup()
	rm := resourcemanager.NewResourceManager(mqttserver.NewMQTTServer(), db)
	return JobController{
		config:          config,
		mailAPI:         mailAPI,
		resourceManager: rm,
		paymentProvider: paypal.Setup(db),
		DataStore:       db,
		member:          member.NewMemberService(db, rm),
	}
}

func (j JobController) CheckMemberSubscriptions() {
	log.Infof("[scheduled-job] checking member subscription status")
	if len(j.config.PaypalURL) == 0 {
		logger.Debug("paypal url isn't set")
		return
	}
	members := j.DataStore.GetMembers()

	for _, member := range members {
		if member.Level == uint8(models.Credited) {
			// do nothing to credited members
			continue
		}
		if member.SubscriptionID == "none" || len(member.SubscriptionID) == 0 {
			// we might need to figure out why they don't have a subscriptionID
			log.Printf("no subscriptionID for: %s", member.Name)
			continue
		}

		status, lastPayment, err := j.paymentProvider.GetSubscription(member.SubscriptionID)
		if err != nil {
			log.Error(err)
			continue
		}

		j.SetMemberLevel(status, lastPayment, member)
	}
}

func (j JobController) SetMemberLevel(status string, lastPayment models.Payment, member models.Member) {
	log.Infof("[scheduled-job] setting member level: %s - %s", member.Name, status)

	switch status {
	case models.ActiveStatus:
		j.member.ActiveStatusHandler(member, lastPayment)
		return
	case models.CanceledStatus:
		j.member.CancelStatusHandler(member, lastPayment)
		return
	case models.SuspendedStatus:
		j.DataStore.SetMemberLevel(member.ID, models.Inactive)
	default:
		return
	}
}

func (j JobController) CheckResourceInit() {
	log.Infof("[scheduled-job] setup mqtt subscriptions to resources")

	resources := j.DataStore.GetResources()

	// on startup we will subscribe to resources and publish an initial status check
	for _, r := range resources {
		j.resourceManager.MQTTServer.Subscribe(r.Name+"/send", j.resourceManager.OnAccessEvent)
		j.resourceManager.MQTTServer.Subscribe(r.Name+"/result", j.resourceManager.HealthCheck)
		j.resourceManager.MQTTServer.Subscribe(r.Name+"/sync", j.resourceManager.OnHeartBeat)
		j.resourceManager.MQTTServer.Subscribe(r.Name+"/cleanup", j.resourceManager.OnRemoveInvalidRequest)
		j.resourceManager.CheckStatus(r)
	}
}

func (j JobController) CheckResourceInterval() {
	log.Infof("[scheduled-job] checking resource status")

	resources := j.DataStore.GetResources()

	for _, r := range resources {
		j.resourceManager.CheckStatus(r)
	}
}

var IPAddressCache string

func (j JobController) CheckIPAddressInterval() {
	log.Infof("[scheduled-job] checking ip address")

	resp, err := http.Get("https://icanhazip.com/")
	if err != nil {
		log.Errorf("can't get IP address: %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	currentIp := strings.TrimSpace(string(body))
	logger.Infof("ip addr: %s", currentIp)

	const ipFileName string = ".public_ip_address"
	// detect if file exists
	_, err = os.Stat(ipFileName)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(ipFileName)
		if err != nil {
			log.Error(err)
			return
		}
		defer file.Close()
	}

	b, err := ioutil.ReadFile(ipFileName)
	if err != nil {
		log.Error(err)
		return
	}

	err = ioutil.WriteFile(ipFileName, body, 0644)
	if err != nil {
		log.Error(err)
		return
	}

	// if this is the first run, don't send an email,
	//   but set the ip address
	previousIp := strings.TrimSpace(string(b))
	if previousIp == "" || previousIp == currentIp {
		return
	}

	ipModel := struct {
		IpAddress string
	}{
		IpAddress: currentIp,
	}

	mailer := mail.NewMailer(j.DataStore, j.mailAPI, j.config)
	mailer.SendCommunication(mail.IpChanged, j.config.AdminEmail, ipModel)
}

func (j JobController) RemovedInvalidUIDs() {
	log.Infof("[scheduled-job] removing any invalid members from resources")
	j.resourceManager.RemovedInvalidUIDs()
}
func (j JobController) EnableValidUIDs() {
	log.Infof("[scheduled-job] enabling valid members on resources")
	j.resourceManager.EnableValidUIDs()
}
func (j JobController) UpdateResources() {
	log.Infof("[scheduled-job] updating resources")
	j.resourceManager.UpdateResources()
}
func (j JobController) UpdateMemberCounts() {
	log.Infof("[scheduled-job] updating member counts")
	j.DataStore.UpdateMemberCounts()
}
