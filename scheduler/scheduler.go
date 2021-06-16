package scheduler

import (
	"io/ioutil"
	"memberserver/config"
	"memberserver/database"
	"memberserver/mail"
	"memberserver/payments"
	"memberserver/resourcemanager"
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

// Setup Scheduler
//  We want certain tasks to happen on a regular basis
//  The scheduler will make sure that happens
func Setup() {
	mailApi, _ = mail.Setup()
	c, _ = config.Load()

	scheduleTask(checkPaymentsInterval*time.Hour, payments.GetPayments, payments.GetPayments)
	scheduleTask(evaluateMemberStatusInterval*time.Hour, checkMemberStatus, checkMemberStatus)
	scheduleTask(resourceStatusCheckInterval*time.Hour, checkResourceInit, checkResourceTick)
	scheduleTask(resourceUpdateInterval*time.Hour, resourcemanager.UpdateResources, resourcemanager.UpdateResources)
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
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
	}
	defer db.Release()

	db.EvaluateMembers()
}

func checkResourceInit() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
		return
	}

	resources := db.GetResources()
	defer db.Release()

	// on startup we will subscribe to resources and publish an initial status check
	for _, r := range resources {
		resourcemanager.Subscribe(r.Name+"/send", resourcemanager.OnAccessEvent)
		resourcemanager.Subscribe(r.Name+"/result", resourcemanager.HealthCheck)
		resourcemanager.Subscribe(r.Name+"/sync", resourcemanager.OnHeartBeat)
		resourcemanager.CheckStatus(r)
	}
}

func checkResourceTick() {
	db, err := database.Setup()
	if err != nil {
		log.Errorf("error setting up db: %s", err)
		return
	}

	resources := db.GetResources()
	defer db.Release()

	for _, r := range resources {
		resourcemanager.CheckStatus(r)
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

	db, err := database.Setup()
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
