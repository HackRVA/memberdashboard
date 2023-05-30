package jobs

import (
	"fmt"
	"io/ioutil"

	config "github.com/HackRVA/memberserver/configs"
	"github.com/HackRVA/memberserver/pkg/membermgr/datastore"
	"github.com/HackRVA/memberserver/pkg/membermgr/integrations"
	"github.com/HackRVA/memberserver/pkg/membermgr/services"
	"github.com/HackRVA/memberserver/pkg/membermgr/services/mail"
	"github.com/HackRVA/memberserver/pkg/paypal"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"

	"net/http"
	"os"
	"strings"
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
		if member.Level == uint8(models.Credited) {
			// do nothing to credited members
			continue
		}
		if member.SubscriptionID == "none" || len(member.SubscriptionID) == 0 {
			howToMsg := "\nUpdate the member's subscriptionID in the webUI.  Or have a db admin update the subscriptionID manually.\n\ne.g."
			exampleQuery := fmt.Sprintf("UPDATE membership.members\n\t\t\tSET subscription_id=<subscriptionID>\n\t\t\tWHERE email='%s';", member.Email)

			j.logger.Errorf("deactivating member (name: %s email: %s) because no subscriptionID was found. \n%s\n\n%s\n\n", member.Name, member.Email, howToMsg, exampleQuery)
			j.SetMemberLevel(models.SuspendedStatus, models.Payment{}, member)
			continue
		}

		status, lastPaymentAmount, lastPaymentTime, err := j.paymentProvider.GetSubscription(member.SubscriptionID)
		if err != nil {
			if member.Level == uint8(models.Inactive) {
				j.logger.Debugf("error getting subscription status for (%s, %s). However, member is already inactive. %s", member.Email, member.Name, err.Error())
				continue
			}
			j.logger.Errorf("error getting subscription: %s (%s, %s) setting to inactive until status is investigated", err.Error(), member.Email, member.Name)
			j.SetMemberLevel(models.SuspendedStatus, models.Payment{}, member)
			continue
		}

		j.SetMemberLevel(status, models.Payment{
			Amount: lastPaymentAmount,
			Time:   lastPaymentTime,
		}, member)
	}
}

func (j JobController) SetMemberLevel(status string, lastPayment models.Payment, member models.Member) {
	j.logger.Infof("[scheduled-job] setting member level: %s - %s - last payment amount: %s", member.Name, status, lastPayment.Amount)

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
