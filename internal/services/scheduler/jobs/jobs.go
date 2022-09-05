package jobs

import (
	"io/ioutil"
	"memberserver/internal/datastore"
	"memberserver/internal/integrations"
	"memberserver/internal/models"
	"memberserver/internal/services/config"
	"memberserver/internal/services/mail"
	"memberserver/internal/services/member"
	"memberserver/internal/services/resourcemanager"
	"memberserver/pkg/mqtt"
	"memberserver/pkg/paypal"
	"memberserver/pkg/slack"

	"net/http"
	"os"
	"strings"
)

type JobController struct {
	config          config.Config
	DataStore       datastore.DataStore
	mailAPI         mail.MailApi
	resourceManager resourcemanager.ResourceManager
	paymentProvider integrations.PaymentProvider
	member          member.MemberService
	logger          logger
}

func New(db datastore.DataStore, logger logger) JobController {
	config, _ := config.Load()
	mailAPI, _ := mail.Setup()
	rm := resourcemanager.New(mqtt.New(), db, slack.Notifier{}, logger)
	pp := paypal.Setup(config.PaypalURL, config.PaypalClientID, config.PaypalClientSecret, logger)
	return JobController{
		config:          config,
		mailAPI:         mailAPI,
		resourceManager: rm,
		paymentProvider: pp,
		DataStore:       db,
		member:          member.New(db, rm, pp, logger),
		logger:          logger,
	}
}

func (j JobController) CheckMemberSubscriptions() {
	j.logger.Infof("[scheduled-job] checking member subscription status")
	if len(j.config.PaypalURL) == 0 {
		j.logger.Debug("paypal url isn't set")
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
			j.logger.Printf("no subscriptionID for: %s", member.Name)
			continue
		}

		status, lastPaymentAmount, lastPaymentTime, err := j.paymentProvider.GetSubscription(member.SubscriptionID)
		if err != nil {
			j.logger.Error(err)
			continue
		}

		j.SetMemberLevel(status, models.Payment{
			Amount: lastPaymentAmount,
			Time:   lastPaymentTime,
		}, member)
	}
}

func (j JobController) SetMemberLevel(status string, lastPayment models.Payment, member models.Member) {
	j.logger.Infof("[scheduled-job] setting member level: %s - %s", member.Name, status)

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
	j.logger.Infof("[scheduled-job] setup mqtt subscriptions to resources")

	resources := j.DataStore.GetResources()

	config, _ := config.Load()

	// on startup we will subscribe to resources and publish an initial status check
	for _, r := range resources {
		j.resourceManager.MQTTServer.Subscribe(config.MQTTBrokerAddress, r.Name+"/send", j.resourceManager.OnAccessEvent)
		j.resourceManager.MQTTServer.Subscribe(config.MQTTBrokerAddress, r.Name+"/result", j.resourceManager.HealthCheck)
		j.resourceManager.MQTTServer.Subscribe(config.MQTTBrokerAddress, r.Name+"/sync", j.resourceManager.OnHeartBeat)
		j.resourceManager.MQTTServer.Subscribe(config.MQTTBrokerAddress, r.Name+"/cleanup", j.resourceManager.OnRemoveInvalidRequest)
		j.resourceManager.CheckStatus(r)
	}
}

func (j JobController) CheckResourceInterval() {
	j.logger.Infof("[scheduled-job] checking resource status")

	resources := j.DataStore.GetResources()

	for _, r := range resources {
		j.resourceManager.CheckStatus(r)
	}
}

var IPAddressCache string

func (j JobController) CheckIPAddressInterval() {
	j.logger.Infof("[scheduled-job] checking ip address")

	resp, err := http.Get("https://icanhazip.com/")
	if err != nil {
		j.logger.Errorf("can't get IP address: %s", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		j.logger.Error(err)
		return
	}

	currentIp := strings.TrimSpace(string(body))
	j.logger.Infof("ip addr: %s", currentIp)

	const ipFileName string = ".public_ip_address"
	// detect if file exists
	_, err = os.Stat(ipFileName)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(ipFileName)
		if err != nil {
			j.logger.Error(err)
			return
		}
		defer file.Close()
	}

	b, err := ioutil.ReadFile(ipFileName)
	if err != nil {
		j.logger.Error(err)
		return
	}

	err = ioutil.WriteFile(ipFileName, body, 0644)
	if err != nil {
		j.logger.Error(err)
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
	j.logger.Infof("[scheduled-job] removing any invalid members from resources")
	j.resourceManager.RemovedInvalidUIDs()
}
func (j JobController) EnableValidUIDs() {
	j.logger.Infof("[scheduled-job] enabling valid members on resources")
	j.resourceManager.EnableValidUIDs()
}
func (j JobController) UpdateResources() {
	j.logger.Infof("[scheduled-job] updating resources")
	j.resourceManager.UpdateResources()
}
func (j JobController) UpdateMemberCounts() {
	j.logger.Infof("[scheduled-job] updating member counts")
	j.DataStore.UpdateMemberCounts()
}
