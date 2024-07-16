package jobs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/pkg/membermgr/datastore"
	"github.com/HackRVA/memberserver/pkg/membermgr/integrations"
	"github.com/HackRVA/memberserver/pkg/membermgr/services"
	"github.com/HackRVA/memberserver/pkg/membermgr/services/mail"
	"github.com/HackRVA/memberserver/pkg/paypal"
)

type JobController struct {
	config          config.Config
	DataStore       datastore.DataStore
	mailAPI         mail.MailApi
	resourceManager services.Resource
	paymentProvider integrations.PaymentProvider
	member          services.Member
	logger          logger
}

func New(db datastore.DataStore, logger logger, member services.Member, resource services.Resource) JobController {
	config, _ := config.Load()
	mailAPI, _ := mail.Setup()
	pp := paypal.Setup(config.PaypalURL, config.PaypalClientID, config.PaypalClientSecret, logger)
	return JobController{
		config:          config,
		mailAPI:         mailAPI,
		resourceManager: resource,
		paymentProvider: pp,
		DataStore:       db,
		member:          member,
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
		j.member.CheckStatus(member.SubscriptionID)
	}
	j.checkActiveMembersWithoutSubscription()
}

func (j JobController) checkActiveMembersWithoutSubscription() {
	membersWithoutSubscription := j.member.GetActiveMembersWithoutSubscription()
	if len(membersWithoutSubscription) == 0 {
		return
	}
	var notification string
	notification += "Found active members without subscription... this needs to be addressed:\n```\n"

	for _, m := range membersWithoutSubscription {
		notification += fmt.Sprintf(
			"ID: %-20s Name: %-20s Email: %-30s RFID: %-15s Level: %-5d SubscriptionID: %-15s\n",
			m.ID, m.Name, m.Email, m.RFID, m.Level, m.SubscriptionID,
		)
	}

	notification += "```\n"

	j.logger.Errorf("%s", notification)
}

func (j JobController) CheckResourceInit() {
	j.logger.Infof("[scheduled-job] setup mqtt subscriptions to resources")

	resources := j.DataStore.GetResources()

	config, _ := config.Load()

	// on startup we will subscribe to resources and publish an initial status check
	for _, r := range resources {
		j.resourceManager.MQTT().Subscribe(config.MQTTBrokerAddress, r.Name+"/send", j.resourceManager.ReceiveHandler)
		j.resourceManager.MQTT().Subscribe(config.MQTTBrokerAddress, r.Name+"/result", j.resourceManager.HealthCheckHandler)
		j.resourceManager.MQTT().Subscribe(config.MQTTBrokerAddress, r.Name+"/sync", j.resourceManager.OnHeartBeatHandler)
		j.resourceManager.MQTT().Subscribe(config.MQTTBrokerAddress, r.Name+"/cleanup", j.resourceManager.OnRemoveInvalidRequestHandler)
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
		file, err := os.Create(ipFileName)
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
