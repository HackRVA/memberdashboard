package scheduler

import (
	"io/ioutil"
	"memberserver/api/models"
	"memberserver/config"
	"memberserver/datastore"
	"memberserver/datastore/dbstore.go"
	"memberserver/datastore/in_memory"
	"memberserver/mail"
	"memberserver/payments"
	"memberserver/resourcemanager"
	"memberserver/resourcemanager/mqttserver"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// checkPaymentsInterval - check the resources every 24 hours
const checkPaymentsInterval = 24

// evaluateMemberStatusInterval - check the resources every 25 hours
const evaluateMemberStatusInterval = 25

// resourceStatusCheckInterval - check the resources every hour
const resourceStatusCheckInterval = 1

const resourceUpdateInterval = 4

// checkIPInterval - check the IP Address daily
const checkIPInterval = 24

var c config.Config
var mailApi mail.MailApi
var db datastore.DataStore
var rm *resourcemanager.ResourceManager

// Setup Scheduler
//  We want certain tasks to happen on a regular basis
//  The scheduler will make sure that happens
func Setup(d datastore.DataStore) {
	mailApi, _ = mail.Setup()
	c, _ = config.Load()
	db = d
	rm = resourcemanager.NewResourceManager(mqttserver.NewMQTTServer(), &in_memory.In_memory{})

	scheduleTask(checkPaymentsInterval*time.Hour, payments.GetPayments, payments.GetPayments)
	scheduleTask(evaluateMemberStatusInterval*time.Hour, checkMemberStatus, checkMemberStatus)
	scheduleTask(resourceStatusCheckInterval*time.Hour, checkResourceInit, checkResourceTick)
	scheduleTask(resourceUpdateInterval*time.Hour, rm.UpdateResources, rm.UpdateResources)
	scheduleTask(checkIPInterval*time.Hour, checkIPAddressTick, checkIPAddressTick)
}

func scheduleTask(interval time.Duration, initFunc func(), tickFunc func()) {
	initFunc()

	// quietly check the resource status on an interval
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				tickFunc()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func checkMemberStatus() {
	db.ApplyMemberCredits()
	db.UpdateMemberTiers()

	const memberGracePeriod = 46
	const membershipMonth = 31

	mailer := mail.NewMailer(db, mailApi, c)

	pendingRevokation, err := db.GetCommunication(mail.PendingRevokationMember.String())

	if err != nil {
		log.Errorf("Unable to get communication %v. Err: %v", mail.PendingRevokationMember, err)
		return
	}

	pastDueAccounts := db.GetPastDueAccounts()
	for _, a := range pastDueAccounts {
		if a.DaysSinceLastPayment > memberGracePeriod {
			mailer.SendCommunication(mail.AccessRevokedLeadership, c.AdminEmail, a)
			mailer.SendCommunication(mail.AccessRevokedMember, a.Email, a)
			db.SetMemberLevel(a.MemberId, models.Inactive)
		} else if a.DaysSinceLastPayment > membershipMonth {
			if !mailer.IsThrottled(pendingRevokation, models.Member{ID: a.MemberId}) {
				//TODO: [ML] Does it make sense to send this to leadership?  It might be like spam...
				mailer.SendCommunication(mail.PendingRevokationLeadership, c.AdminEmail, a)
				mailer.SendCommunication(mail.PendingRevokationMember, a.Email, a)
			}
		}
	}
}

func checkResourceInit() {
	resources := db.GetResources()

	// on startup we will subscribe to resources and publish an initial status check
	for _, r := range resources {
		rm.MQTTServer.Subscribe(r.Name+"/send", resourcemanager.OnAccessEvent)
		rm.MQTTServer.Subscribe(r.Name+"/result", resourcemanager.HealthCheck)
		rm.MQTTServer.Subscribe(r.Name+"/sync", resourcemanager.OnHeartBeat)
		rm.CheckStatus(r)
	}
}

func checkResourceTick() {
	resources := db.GetResources()

	for _, r := range resources {
		rm.CheckStatus(r)
	}
}

var IPAddressCache string

func checkIPAddressTick() {
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
	log.Debugf("ip addr: %s", currentIp)

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

	db, err := dbstore.Setup()
	if err != nil {
		log.Printf("Err: %v", err)
	}

	ipModel := struct {
		IpAddress string
	}{
		IpAddress: currentIp,
	}

	mailer := mail.NewMailer(db, mailApi, c)
	mailer.SendCommunication(mail.IpChanged, c.AdminEmail, ipModel)
}
